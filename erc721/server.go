package erc721

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
)

// A MetadataServer handles HTTP routes to serve ERC721 metadata JSON and
// associated images. If a contract binding is provided, it is checked to ensure
// that the requested token already exists, thus allowing a MetadataServer to be
// used for delayed reveals.
type MetadataServer struct {
	// BaseURL is the base URL of the server; i.e. everything except the path,
	// which will be overwritten.
	BaseURL *url.URL
	// MetadataPath and ImagePath are the routing paths for metadata and images
	// respectively. They follow syntax of github.com/julienschmidt/httprouter
	// and use the TokenIDParam to extract the token ID.
	MetadataPath, ImagePath string
	// TokenIDBase, if non-zero, uses a custom base for decoding the token ID,
	// defaulting to base 10.
	TokenIDBase int
	// The Contract, if provided, is used to confirm that tokens exist before
	// responding with metadata or images. Checks use the ownerOf function,
	// which must not revert.
	Contract Interface

	// Metadata and Image are responsible for returning a token's metadata and
	// image, respectively (surprise, surprise!). If Contract is non-nil, the
	// token is guaranteed to exist if Metadata/Image is called. Only 200, 400,
	// 404 and 500 are allowed as HTTP codes, and these will be propagated to
	// the end user.
	Metadata func(Interface, *TokenID, httprouter.Params) (md *Metadata, httpCode int, err error)
	Image    func(Interface, *TokenID, httprouter.Params) (img io.Reader, contentType string, httpCode int, err error)
}

// ListenAndServe returns http.ListenAndServe(addr, s.Handler()).
func (s *MetadataServer) ListenAndServe(addr string) error {
	h, err := s.Handler()
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, h)
}

// Handler returns a Handler, for use with http.ListenAndServe(), that handles
// all requests for metadata and images. Unless the Handler is specifically
// needed for non-default uses, prefer s.ListenAndServer().
func (s *MetadataServer) Handler() (http.Handler, error) {
	r := httprouter.New()

	for name, path := range map[string]string{
		"Metadata": s.MetadataPath,
		"Image":    s.ImagePath,
	} {
		if !strings.Contains(path, fullTokenIDParam) {
			return nil, fmt.Errorf("%sPath %q must contain %q", name, path, fullTokenIDParam)
		}
	}

	r.GET(s.MetadataPath, httpErrHandler(s.metadata))
	r.GET(s.ImagePath, httpErrHandler(s.images))

	return r, nil
}

// httpErrHandler allows httprouter.Handle-like functions to return errors. If
// the returned error is of the type *httpError then its code is propagated; 400
// and 404 errors also have their message propagated to the client. All other
// codes are hashed and logged, with only a portion of the hash returned to the
// end user.
func httpErrHandler(fn func(http.ResponseWriter, *http.Request, httprouter.Params) error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		obfuscate := func(code int, msg string) {
			id := crypto.Keccak256([]byte(msg))
			id = id[:8]
			glog.Errorf("%x: %s", id, msg)
			http.Error(w, fmt.Sprintf("see log: %x", id), code)
		}

		switch err := fn(w, r, p).(type) {
		case nil:
		case *httpError:
			switch err.code {
			case http.StatusBadRequest, http.StatusNotFound:
				http.Error(w, err.msg, err.code)
			default:
				obfuscate(err.code, err.msg)
			}
		default:
			obfuscate(500, err.Error())
		}
	}
}

// httpError is an error that carries an HTTP response code and a message.
type httpError struct {
	code int
	msg  string
}

func (e *httpError) Error() string {
	return e.msg
}

// errorf returns an httpError if code != 200, otherwise it returns nil.
func errorf(code int, format string, a ...interface{}) error {
	if code == 200 {
		return nil
	}
	return &httpError{
		code: code,
		msg:  fmt.Sprintf(format, a...),
	}
}

// A tokenDataFunc returns arbitrary HTTP response data for a token.
type tokenDataFunc func(Interface, *TokenID, httprouter.Params) (body io.Reader, contentType string, code int, err error)

// tokenDataHandler is a generic handler for any token data, abstracting shared
// logic from the metadata and image handlers.
func (s *MetadataServer) tokenDataHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, fnName string, fn tokenDataFunc) error {
	id, err := s.tokenID(params)
	if err != nil {
		return errorf(500, "%T.tokenID(%+v): %v", s, params, err)
	}
	if id == nil {
		return errorf(404, "token %q not minted", params.ByName(TokenIDParam))
	}

	if fn == nil {
		return errorf(400, "unsupported method %s", fnName)
	}

	body, contentType, code, err := fn(s.Contract, id, params)
	if err != nil {
		return errorf(500, "%s(%s): %v", fnName, id, err)
	}

	switch code {
	case 200:
	case 400, 404, 500:
		return errorf(code, "%s", err)
	default:
		return errorf(500, "unsupported code %d returned by %s(%s)", code, fnName, id)
	}

	w.Header().Add("Content-Type", contentType)
	if _, err := io.Copy(w, body); err != nil {
		return errorf(500, "io.Copy([http response], [%s data]): %v", fnName, err)
	}
	return nil
}

// metadata handles requests for metadata, sourcing it from the user-provided
// s.Metadata() function, and substituting the Image field appropriately such
// that it will point to the MetadataServer's image endpoint.
func (s *MetadataServer) metadata(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	return s.tokenDataHandler(w, r, params, "Metadata", func(i Interface, id *TokenID, params httprouter.Params) (io.Reader, string, int, error) {
		md, code, err := s.Metadata(s.Contract, id, params)
		if err != nil {
			return nil, "", code, err
		}

		if md.Image == "" {
			img := *s.BaseURL
			img.Path = strings.ReplaceAll(s.ImagePath, fullTokenIDParam, id.Text(s.tokenIDBase()))
			md.Image = img.String()
		}

		buf, err := json.Marshal(md)
		if err != nil {
			return nil, "", 500, fmt.Errorf("json.Marshal(%T = %+v): %v", md, md, err)
		}
		return bytes.NewReader(buf), "application/json", 200, nil
	})
}

// images handles requests for images, sourcing them from the user-provided
// s.Images() function.
func (s *MetadataServer) images(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	return s.tokenDataHandler(w, r, params, "Image", s.Image)
}

// TokenIDParam is the name of the httprouter parameter matched by metadata and
// image endpoints. Examples of valid paths include:
//
//  /metadata/:tokenId
//  /images/:tokenId
//  /path/to/metadata/:tokenId/:otherParam/passed/to/user/functions
const TokenIDParam = "tokenId"

const fullTokenIDParam = ":" + TokenIDParam

// tokenID extracts the `TokenIDParam` from the params. If s.Contract is
// non-nil, it is used to check that the token already exists—if not then
// tokenID() returns (nil, nil).
func (s *MetadataServer) tokenID(params httprouter.Params) (*TokenID, error) {
	rawID := params.ByName(TokenIDParam)
	if rawID == "" {
		return nil, fmt.Errorf("no %q param", TokenIDParam)
	}

	base := s.tokenIDBase()
	id, ok := new(big.Int).SetString(rawID, base)
	if !ok {
		return nil, fmt.Errorf("token ID %q not parsed in base %d", rawID, base)
	}

	if s.Contract != nil {
		if _, err := s.Contract.OwnerOf(nil, id); err != nil {
			return nil, nil
		}
	}
	return TokenIDFromBig(id)
}

// tokenIDBase returns s.TokenIDBase if non-zero, otherwise it returns 10.
func (s *MetadataServer) tokenIDBase() int {
	switch b := s.TokenIDBase; b {
	case 0:
		return 10
	default:
		return b
	}
}
