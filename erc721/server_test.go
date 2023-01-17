package erc721

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/divergencetech/ethier/ethtest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/julienschmidt/httprouter"

	contract "github.com/divergencetech/ethier/tests/erc721"
)

func deploy(t *testing.T, totalSupply int64) Interface {
	t.Helper()

	sim := ethtest.NewSimulatedBackendTB(t, 1)
	_, _, nft, err := contract.DeployTestableERC721ACommon(sim.Acc(0), sim, common.Address{1}, big.NewInt(0))
	if err != nil {
		t.Fatalf("DeployTestableERC721ACommon(): %v", err)
	}

	sim.Must(t, "%T.MintN(%d)", nft, totalSupply)(nft.MintN(sim.Acc(0), big.NewInt(totalSupply)))
	return nft
}

// start starts a new http server with requests handled by srv.Handler(), and
// returns the base URL of the started server.
func start(t *testing.T, srv *Server) string {
	handler, err := srv.Handler()
	if err != nil {
		t.Fatalf("%T{%+v}.Handler() error %v", srv, srv, err)
	}

	httpSrv := httptest.NewServer(handler)
	t.Cleanup(httpSrv.Close)

	base, err := url.Parse(httpSrv.URL)
	if err != nil {
		t.Fatalf("url.Parse(%q = %T.URL): %v", httpSrv.URL, httpSrv, err)
	}
	srv.BaseURL = base

	return httpSrv.URL
}

// httpGet calls http.Get(url), reports any errors on t.Fatal(), and returns the
// HTTP response.
func httpGet(t *testing.T, url string) *http.Response {
	t.Helper()

	res, err := http.Get(url)
	if err != nil {
		t.Fatalf("http.Get(%q): %v", url, err)
	}
	return res
}

// testContentType asserts that the response has 200 code and the expected
// content-type header.
func testContentType(t *testing.T, resp *http.Response, wantContentType string) {
	t.Helper()

	if got, want := resp.StatusCode, 200; got != want {
		t.Errorf("HTTP GET %q; got code %d; want %d", resp.Request.URL, got, want)
	}
	if got, want := resp.Header.Get("Content-Type"), wantContentType; got != want {
		t.Errorf("HTTP GET %q; got Content-Type %q; want %q", resp.Request.URL, got, want)
	}
}

// metadataFromResponse parses the http.Response.Body into a new Metadata
// instance and returns it.
func metadataFromResponse(t *testing.T, r *http.Response) *Metadata {
	t.Helper()

	dec := json.NewDecoder(r.Body)
	md := new(Metadata)
	if err := dec.Decode(md); err != nil {
		t.Fatalf("%T.Decode(%T) error %v", dec, md, err)
	}
	return md
}

func TestMetadataServer(t *testing.T) {
	const (
		totalSupply = 16
		imageType   = "image/png"
	)

	tests := []struct {
		Name                  string
		ExternalImageURL      string
		Image                 ImageHandler
		wantInternalImageCode int
	}{
		{
			Name: "Internal images",
			Image: func(_ Interface, id *TokenID, params httprouter.Params) (io.Reader, string, int, error) {
				return strings.NewReader(fmt.Sprintf("Image %s", id)), imageType, 200, nil
			},
			wantInternalImageCode: 200,
		},
		{
			Name:                  "External images",
			ExternalImageURL:      "http://foo.bar",
			wantInternalImageCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			srv := &Server{
				BaseURL:     nil, // will be set to the test-server URL by start()
				TokenIDBase: 16,
				Contract:    deploy(t, totalSupply),
				Metadata: []MetadataEndpoint{{
					Path: "/metadata/:tokenId",
					Handler: func(_ Interface, id *TokenID, params httprouter.Params) (*Metadata, int, error) {
						md := Metadata{
							Name:  fmt.Sprintf("Token %s", id),
							Image: tt.ExternalImageURL,
						}

						return &md, 200, nil
					},
				}},
				Image: []ImageEndpoint{{
					Path:    "/image/:tokenId",
					Handler: tt.Image,
				}},
			}
			baseURL := start(t, srv)
			tokenURL := func(id int) string {
				// Note the use of %x as we set TokenIDBase to 16.
				return fmt.Sprintf("%s/metadata/%x", baseURL, id)
			}
			internalImageURL := func(id int) string {
				// Note the use of %x as we set TokenIDBase to 16.
				return fmt.Sprintf("%s/image/%x", baseURL, id)
			}

			for id := 0; id < totalSupply; id++ {
				t.Run(fmt.Sprintf("token %d", id), func(t *testing.T) {
					// The image path has to be extracted from the metadata response,
					// but it's cleaner to separate into sub tests and we therefore need
					// a variable at a higher scope.
					var gotMetadata *Metadata

					t.Run("metadata", func(t *testing.T) {
						path := tokenURL(id)
						resp := httpGet(t, path)
						testContentType(t, resp, "application/json")

						gotMetadata = metadataFromResponse(t, resp)
						want := &Metadata{
							Name:  fmt.Sprintf("Token %d", id),
							Image: tt.ExternalImageURL,
						}

						var opts []cmp.Option
						if tt.ExternalImageURL == "" {
							// We ignore the image field when using the internal
							// endpoint because we test its validity later.
							opts = append(opts, cmpopts.IgnoreFields(Metadata{}, "Image"))
						}

						if diff := cmp.Diff(want, gotMetadata, opts...); diff != "" {
							t.Errorf("HTTP GET %q; parsed %T diff (-want +got):\n%s", path, gotMetadata, diff)
						}
					})

					if t.Failed() {
						// Running image tests without correct Metadata is guaranteed to
						// fail, so don't spuriously pollute the error messages.
						return
					}

					t.Run("image", func(t *testing.T) {
						if tt.ExternalImageURL != "" {
							// Skiping this test for external images because we
							// already checked if the correct URL is returned.
							return
						}

						resp := httpGet(t, gotMetadata.Image)
						testContentType(t, resp, imageType)

						got, err := io.ReadAll(resp.Body)
						if err != nil {
							t.Fatalf("io.ReadAll([http response body]): %v", err)
						}
						if want := fmt.Sprintf("Image %d", id); string(got) != want {
							t.Errorf("HTTP GET %q; got body %q; want %q", gotMetadata.Image, got, want)
						}

						t.Run("path", func(t *testing.T) {
							u, err := url.Parse(gotMetadata.Image)
							if err != nil {
								t.Fatalf("url.Parse(%q = %T.Image return value): %v", gotMetadata.Image, gotMetadata, err)
							}
							if got, want := u.Path, fmt.Sprintf("/image/%x", id); got != want {
								t.Errorf("%T.Image path = %q; want %q", gotMetadata, got, want)
							}
						})
					})

					t.Run("internal image code", func(t *testing.T) {
						path := internalImageURL(id)
						resp := httpGet(t, path)
						if got, want := resp.StatusCode, tt.wantInternalImageCode; got != want {
							t.Errorf("HTTP GET %q: got code %v, want %v", path, got, want)
						}
					})

				})
			}

			t.Run("non-existent token", func(t *testing.T) {
				t.Run("metadata", func(t *testing.T) {
					path := tokenURL(totalSupply)
					resp := httpGet(t, path)
					if got, want := resp.StatusCode, 404; got != want {
						t.Errorf("HTTP GET %q (non-existent token) got code %d; want %d", path, got, want)
					}
				})
			})

		})
	}
}

func TestMultipleMetadataEndpoints(t *testing.T) {
	srv := &Server{
		Metadata: []MetadataEndpoint{
			{
				Path: "/default/:tokenId",
				Handler: func(_ Interface, id *TokenID, _ httprouter.Params) (*Metadata, int, error) {
					md := Metadata{
						Name: fmt.Sprintf("DEFAULT %s", id),
					}
					return &md, 200, nil
				},
			},
			{
				Path: "/extra/:tokenId",
				Handler: func(_ Interface, id *TokenID, _ httprouter.Params) (*Metadata, int, error) {
					md := Metadata{
						Name: fmt.Sprintf("EXTRA %s", id),
					}
					return &md, 200, nil
				},
			},
			{
				Path: "/with-explicit-image/:tokenId",
				Handler: func(_ Interface, id *TokenID, _ httprouter.Params) (*Metadata, int, error) {
					md := Metadata{
						Name:  fmt.Sprintf("EXPLICIT %s", id),
						Image: fmt.Sprintf("explicit-image-path/%s", id),
					}
					return &md, 200, nil
				},
			},
		},
		Image: []ImageEndpoint{
			{
				// This first one MUST be added to the Metadata if the Image
				// field is empty.
				Path: "/images/:tokenId",
			},
			{
				Path: "/other-images/:tokenId",
			},
		},
	}

	baseURL := start(t, srv)
	imagePath := func(id int) string {
		return fmt.Sprintf("%s/images/%d", baseURL, id)
	}

	ignore := cmpopts.IgnoreFields(Metadata{})
	tests := []struct {
		path string
		want *Metadata
	}{
		{
			path: "default/0",
			want: &Metadata{
				Name:  "DEFAULT 0",
				Image: imagePath(0),
			},
		},
		{
			path: "default/42",
			want: &Metadata{
				Name:  "DEFAULT 42",
				Image: imagePath(42),
			},
		},
		{
			path: "extra/42",
			want: &Metadata{
				Name:  "EXTRA 42",
				Image: imagePath(42),
			},
		},
		{
			path: "with-explicit-image/99",
			want: &Metadata{
				Name:  "EXPLICIT 99",
				Image: "explicit-image-path/99",
			},
		},
	}

	for _, tt := range tests {
		url := fmt.Sprintf("%s/%s", baseURL, tt.path)
		got := metadataFromResponse(t, httpGet(t, url))

		if diff := cmp.Diff(tt.want, got, ignore); diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}

func TestInternalImage(t *testing.T) {
	tests := []struct {
		name            string
		metadataHandler MetadataHandler
		wantImage       func(baseURL string, id int) string
	}{
		{
			name: "Empty",
			metadataHandler: func(_ Interface, id *TokenID, _ httprouter.Params) (*Metadata, int, error) {
				md := Metadata{
					Name: fmt.Sprintf("DEFAULT %s", id),
				}
				return &md, 200, nil
			},
			wantImage: func(baseURL string, id int) string {
				return fmt.Sprintf("%s/images/%d", baseURL, id)
			},
		},
		{
			name: "With animation",
			metadataHandler: func(_ Interface, id *TokenID, _ httprouter.Params) (*Metadata, int, error) {
				md := Metadata{
					Name:         fmt.Sprintf("DEFAULT %s", id),
					AnimationURL: fmt.Sprintf("foo://bar/%s", id),
				}
				return &md, 200, nil
			},
			wantImage: func(baseURL string, id int) string {
				return ""
			},
		},
		{
			name: "With image",
			metadataHandler: func(_ Interface, id *TokenID, _ httprouter.Params) (*Metadata, int, error) {
				md := Metadata{
					Name:  fmt.Sprintf("DEFAULT %s", id),
					Image: fmt.Sprintf("foo://bar/%s", id),
				}
				return &md, 200, nil
			},
			wantImage: func(baseURL string, id int) string {
				return fmt.Sprintf("foo://bar/%d", id)
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			srv := &Server{
				Metadata: []MetadataEndpoint{
					{
						Path:    "/metadata/:tokenId",
						Handler: tt.metadataHandler,
					},
				},
				Image: []ImageEndpoint{
					{
						Path: "/images/:tokenId",
					},
				},
			}
			baseURL := start(t, srv)

			for id := 0; id < 20; id++ {
				t.Run(fmt.Sprintf("token %d", id), func(t *testing.T) {
					url := fmt.Sprintf("%s/metadata/%d", baseURL, id)
					got := metadataFromResponse(t, httpGet(t, url))

					if diff := cmp.Diff(tt.wantImage(baseURL, id), got.Image); diff != "" {
						t.Errorf("(-want +got):\n%s", diff)
					}
				})
			}
		})
	}
}
