package authconfig

import (
	"reflect"
	"testing"

	"gopkg.in/Knetic/govaluate.v2"

	"github.com/freifunkMUC/wg-access-server/pkg/authnz/authsession"
)

func Test_evaluateClaimMapping(t *testing.T) {
	type args struct {
		claimMapping map[string]ruleExpression
		oidcClaims   map[string]interface{}
	}

	expr, _ := govaluate.NewEvaluableExpression("'WireguardAdmins' in group_membership")

	tests := []struct {
		name    string
		args    args
		want    authsession.Claims
		wantErr bool
	}{
		{
			args: args{
				claimMapping: map[string]ruleExpression{"admin": {expr}},
				oidcClaims:   map[string]interface{}{"group_membership": []interface{}{"wgas", "WireguardAdmins"}},
			},
			want: authsession.Claims{{Name: "admin", Value: "true"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateClaimMapping(tt.args.claimMapping, tt.args.oidcClaims)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateClaimMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("evaluateClaimMapping() got = %v, want %v", got, tt.want)
			}
		})
	}
}
