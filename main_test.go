package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_replaceVersions(t *testing.T) {
	type args struct {
		in          []byte
		placeholder string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			args: args{
				in: []byte(`FROM caddy:2.4.0-builder AS builder

RUN xcaddy build \
	--with github.com/greenpau/caddy-auth-jwt \
	--with github.com/greenpau/caddy-auth-portal \
	--with github.com/caddy-dns/cloudflare

FROM caddy:2.3.0

COPY --from=builder /usr/bin/caddy /usr/bin/caddy`),
				placeholder: "NA",
			},
			name: "sync to builder",
			want: []byte(`FROM caddy:2.4.0-builder AS builder

RUN xcaddy build \
	--with github.com/greenpau/caddy-auth-jwt \
	--with github.com/greenpau/caddy-auth-portal \
	--with github.com/caddy-dns/cloudflare

FROM caddy:2.4.0

COPY --from=builder /usr/bin/caddy /usr/bin/caddy`),
		},
		{
			args: args{
				in: []byte(`FROM caddy:2.2.0-builder AS builder

RUN xcaddy build \
	--with github.com/greenpau/caddy-auth-jwt \
	--with github.com/greenpau/caddy-auth-portal \
	--with github.com/caddy-dns/cloudflare

FROM caddy:2.3.0

COPY --from=builder /usr/bin/caddy /usr/bin/caddy`),
				placeholder: "NA",
			},
			name: "sync to caddy",
			want: []byte(`FROM caddy:2.3.0-builder AS builder

RUN xcaddy build \
	--with github.com/greenpau/caddy-auth-jwt \
	--with github.com/greenpau/caddy-auth-portal \
	--with github.com/caddy-dns/cloudflare

FROM caddy:2.3.0

COPY --from=builder /usr/bin/caddy /usr/bin/caddy`),
		},
		{
			args: args{
				in: []byte(`FROM caddy:2.3.0-builder AS builder

RUN xcaddy build \
	--with github.com/greenpau/caddy-auth-jwt \
	--with github.com/greenpau/caddy-auth-portal \
	--with github.com/caddy-dns/cloudflare

FROM caddy:2.3.0

COPY --from=builder /usr/bin/caddy /usr/bin/caddy`),
				placeholder: "NA",
			},
			name: "already sync",
			want: []byte(`FROM caddy:2.3.0-builder AS builder

RUN xcaddy build \
	--with github.com/greenpau/caddy-auth-jwt \
	--with github.com/greenpau/caddy-auth-portal \
	--with github.com/caddy-dns/cloudflare

FROM caddy:2.3.0

COPY --from=builder /usr/bin/caddy /usr/bin/caddy`),
		},
		{
			args: args{
				in: []byte(`FROM caddy:2.3.0-builder AS builder

RUN xcaddy build \
	--with github.com/greenpau/caddy-auth-jwt \
	--with github.com/greenpau/caddy-auth-portal \
	--with github.com/caddy-dns/cloudflare

FROM caddy:VERSION

COPY --from=builder /usr/bin/caddy /usr/bin/caddy`),
				placeholder: "VERSION",
			},
			name: "placeholder",
			want: []byte(`FROM caddy:2.3.0-builder AS builder

RUN xcaddy build \
	--with github.com/greenpau/caddy-auth-jwt \
	--with github.com/greenpau/caddy-auth-portal \
	--with github.com/caddy-dns/cloudflare

FROM caddy:2.3.0

COPY --from=builder /usr/bin/caddy /usr/bin/caddy`),
		},
		{
			args: args{
				in:          []byte(`nothing interesting`),
				placeholder: "NA",
			},
			name: "nothing interesting",
			want: []byte(`nothing interesting`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceVersions(tt.args.in, tt.args.placeholder)
			if diff := cmp.Diff(string(tt.want), string(got)); diff != "" {
				t.Errorf("replaceVersions() returns expected result, (-want +got):\n%s", diff)
			}
		})
	}
}
