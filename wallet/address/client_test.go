package address

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/bithubhq/bithub-go/coin"
	"github.com/bithubhq/bithub-go/pkg/testutil"
)

const baseURL = "https://api.coin.com/wallet"

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestClient_CreateAddress(t *testing.T) {
	const expectedEndpoint = "https://api.coin.com/wallet/v1/addresses"

	type fields struct {
		Config *ClientConfig
	}

	type args struct {
		AddressParams AddressParams
	}

	noLabelClient := NewTestClient(func(req *http.Request) *http.Response {
		testutil.Equals(t, expectedEndpoint, req.URL.String())
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"label": "", "address": "foobar"}`)),
			Header:     make(http.Header),
		}
	})

	withLabelClient := NewTestClient(func(req *http.Request) *http.Response {
		testutil.Equals(t, expectedEndpoint, req.URL.String())
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"label": "foobarbaz", "address": "foobar"}`)),
			Header:     make(http.Header),
		}
	})

	errClient := NewTestClient(func(req *http.Request) *http.Response {
		testutil.Equals(t, expectedEndpoint, req.URL.String())
		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error": "some foobar error"}`)),
			Header:     make(http.Header),
		}
	})

	tests := []struct {
		name        string
		args        args
		fields      fields
		wantAddress *Address
		wantErr     bool
	}{
		{
			name: "returns address without label",
			fields: fields{Config: &ClientConfig{
				BaseURL:    baseURL,
				HTTPClient: noLabelClient,
			}},
			args: args{AddressParams: AddressParams{
				Coin: coin.Bitcoin,
			}},
			wantAddress: &Address{
				Label:   "",
				Address: "foobar",
			},
			wantErr: false,
		},
		{
			name: "returns address with label",
			fields: fields{Config: &ClientConfig{
				BaseURL:    baseURL,
				HTTPClient: withLabelClient,
			}},
			args: args{AddressParams: AddressParams{
				Coin:  coin.Bitcoin,
				Label: "foobarbaz",
			}},
			wantAddress: &Address{
				Label:   "foobarbaz",
				Address: "foobar",
			},
			wantErr: false,
		},
		{
			name: "returns error",
			fields: fields{Config: &ClientConfig{
				BaseURL:    baseURL,
				HTTPClient: errClient,
			}},
			wantAddress: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Config: tt.fields.Config,
			}
			gotAddress, err := c.Create(tt.args.AddressParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAddress, tt.wantAddress) {
				t.Errorf("Create() gotAddress = %+v, want %+v", gotAddress, tt.wantAddress)
			}
		})
	}
}
