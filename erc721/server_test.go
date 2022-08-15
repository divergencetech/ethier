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

func start(t *testing.T, srv *MetadataServer) (string, func(*testing.T, string) *http.Response) {
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

	return httpSrv.URL, func(t *testing.T, url string) *http.Response {
		t.Helper()

		res, err := http.Get(url)
		if err != nil {
			t.Fatalf("http.Get(%q): %v", url, err)
		}
		return res
	}
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

func TestMetadataServer(t *testing.T) {
	const (
		totalSupply = 16
		imageType   = "image/png"
	)

	tests := []struct {
		Name             string
		ExternalImageURL string
	}{
		{
			Name: "Internal images",
		},
		{
			Name:             "External images",
			ExternalImageURL: "http://foo.bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			srv := &MetadataServer{
				BaseURL:      nil, // will be set to the test-server URL by start()
				MetadataPath: "/metadata/:tokenId",
				ImagePath:    "/image/:tokenId",
				TokenIDBase:  16,
				Contract:     deploy(t, totalSupply),
				Metadata: func(_ Interface, id *TokenID, params httprouter.Params) (*Metadata, int, error) {
					md := Metadata{
						Name:  fmt.Sprintf("Token %s", id),
						Image: tt.ExternalImageURL,
					}

					return &md, 200, nil
				},
				Image: func(_ Interface, id *TokenID, params httprouter.Params) (io.Reader, string, int, error) {
					return strings.NewReader(fmt.Sprintf("Image %s", id)), imageType, 200, nil
				},
			}
			baseURL, get := start(t, srv)
			tokenURL := func(id int) string {
				// Note the use of %x as we set TokenIDBase to 16.
				return fmt.Sprintf("%s/metadata/%x", baseURL, id)
			}

			for id := 0; id < totalSupply; id++ {
				t.Run(fmt.Sprintf("token %d", id), func(t *testing.T) {
					// The image path has to be extracted from the metadata response,
					// but it's cleaner to separate into sub tests and we therefore need
					// a variable at a higher scope.
					var gotMetadata *Metadata

					t.Run("metadata", func(t *testing.T) {
						path := tokenURL(id)
						resp := get(t, path)
						testContentType(t, resp, "application/json")

						buf, err := io.ReadAll(resp.Body)
						if err != nil {
							t.Fatalf("io.ReadAll([http response body]): %v", err)
						}
						t.Logf("HTTP response body:\n%s", string(buf))
						gotMetadata = new(Metadata)
						if err := json.Unmarshal(buf, gotMetadata); err != nil {
							t.Fatalf("json.Unmarshal([http response body], %T): %v", gotMetadata, err)
						}

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

						resp := get(t, gotMetadata.Image)
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
				})
			}

			t.Run("non-existent token", func(t *testing.T) {
				t.Run("metadata", func(t *testing.T) {
					path := tokenURL(totalSupply)
					resp := get(t, path)
					if got, want := resp.StatusCode, 404; got != want {
						t.Errorf("HTTP GET %q (non-existent token) got code %d; want %d", path, got, want)
					}
				})
			})

		})
	}
}
