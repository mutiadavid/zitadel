package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	zerrors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type OIDCClient struct {
	InstanceID               string                     `json:"instance_id,omitempty"`
	AppID                    string                     `json:"app_id,omitempty"`
	State                    domain.AppState            `json:"state,omitempty"`
	ClientID                 string                     `json:"client_id,omitempty"`
	ClientSecret             *crypto.CryptoValue        `json:"client_secret,omitempty"`
	RedirectURIs             []string                   `json:"redirect_uris,omitempty"`
	ResponseTypes            []domain.OIDCResponseType  `json:"response_types,omitempty"`
	GrantTypes               []domain.OIDCGrantType     `json:"grant_types,omitempty"`
	ApplicationType          domain.OIDCApplicationType `json:"application_type,omitempty"`
	AuthMethodType           domain.OIDCAuthMethodType  `json:"auth_method_type,omitempty"`
	PostLogoutRedirectURIs   []string                   `json:"post_logout_redirect_uris,omitempty"`
	IsDevMode                bool                       `json:"is_dev_mode,omitempty"`
	AccessTokenType          domain.OIDCTokenType       `json:"access_token_type,omitempty"`
	AccessTokenRoleAssertion bool                       `json:"access_token_role_assertion,omitempty"`
	IDTokenRoleAssertion     bool                       `json:"id_token_role_assertion,omitempty"`
	IDTokenUserinfoAssertion bool                       `json:"id_token_userinfo_assertion,omitempty"`
	ClockSkew                time.Duration              `json:"clock_skew,omitempty"`
	AdditionalOrigins        []string                   `json:"additional_origins,omitempty"`
	PublicKeys               map[string][]byte          `json:"public_keys,omitempty"`
	ProjectID                string                     `json:"project_id,omitempty"`
	ProjectRoleKeys          []string                   `json:"project_role_keys,omitempty"`
	Settings                 *OIDCSettings              `json:"settings,omitempty"`
}

//go:embed embed/oidc_client_by_id.sql
var oidcClientQuery string

func (q *Queries) GetOIDCClientByID(ctx context.Context, clientID string, getKeys bool) (client *OIDCClient, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	client, err = database.QueryJSONObject[OIDCClient](ctx, q.client, oidcClientQuery,
		authz.GetInstance(ctx).InstanceID(), clientID, getKeys,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, zerrors.ThrowNotFound(err, "QUERY-wu6Ee", "Errors.App.NotFound")
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-ieR7R", "Errors.Internal")
	}
	return client, err
}
