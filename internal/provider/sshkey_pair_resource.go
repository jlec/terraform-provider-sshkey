package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jlec/terraform-provider-sshkey/internal/keygen"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SSHKeyPairResource{}
	_ resource.ResourceWithImportState = &SSHKeyPairResource{}
)

const (
	sshKeyPath string = "/dev/null"
	sshKeyPass string = ""
)

func NewSSHKeyPairResource() resource.Resource { //nolint:ireturn
	return &SSHKeyPairResource{}
}

// SSHKeyPairResource defines the resource implementation.
type SSHKeyPairResource struct{}

// SSHKeyPairResourceModel describes the resource data model.
type SSHKeyPairResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	Bits              types.Int64  `tfsdk:"bits"`
	Comment           types.String `tfsdk:"comment"`
	PrivateKeyPEM     types.String `tfsdk:"private_key"`
	PublicKey         types.String `tfsdk:"public_key"`
	FingerprintMD5    types.String `tfsdk:"fingerprint_md5"`
	FingerprintSHA256 types.String `tfsdk:"fingerprint_sha256"`
}

func (r *SSHKeyPairResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_pair"
}

//
//nolint:funlen
func (r *SSHKeyPairResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Openssh key resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSHKey identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"comment": schema.StringAttribute{
				Description:         "SSH key comment",
				MarkdownDescription: "SSH key comment",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description:         "SSH key type",
				MarkdownDescription: "SSH key type. Supported types are `rsa` and `ed25519`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(keygen.SSHKeyTypesStrings...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bits": schema.Int64Attribute{
				MarkdownDescription: "When `algorithm` is `RSA`, the size of the generated RSA key, in bits (default: `2048`).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(keygen.SSHRsaBits...),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"private_key": schema.StringAttribute{
				Description:         "OpenSSH private key",
				MarkdownDescription: "OpenSSH private key",
				Computed:            true,
				Sensitive:           true,
			},
			"public_key": schema.StringAttribute{
				Description:         "OpenSSH public key",
				MarkdownDescription: "OpenSSH public key",
				Computed:            true,
			},
			"fingerprint_md5": schema.StringAttribute{
				Description:         "OpenSSH key md5 fingerprint",
				MarkdownDescription: "OpenSSH key md5 fingerprint",
				Computed:            true,
			},
			"fingerprint_sha256": schema.StringAttribute{
				Description:         "OpenSSH key sha256 fingerprint",
				MarkdownDescription: "OpenSSH key sha256 fingerprint",
				Computed:            true,
			},
		},
	}
}

func (r *SSHKeyPairResource) Configure(
	_ context.Context,
	_ resource.ConfigureRequest,
	_ *resource.ConfigureResponse,
) {
}

//
//nolint:funlen
func (r *SSHKeyPairResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var (
		err    error
		data   *SSHKeyPairResourceModel
		sshkey *keygen.SSHKeyPair
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ktyp keygen.KeyType

	switch data.Type.ValueString() {
	case "rsa":
		ktyp = keygen.RSA
	case "ed25519":
		ktyp = keygen.ED25519
	case "ecdsa":
		ktyp = keygen.ECDSA
	}

	if data.Bits.IsUnknown() {
		data.Bits = types.Int64Value(keygen.RsaDefaultBits)
	}

	conf := keygen.SSHKeyPairConfig{
		Passphrase: []byte(sshKeyPass),
		Type:       ktyp,
		Bits:       uint16(data.Bits.ValueInt64()),
	}

	if data.Comment.IsNull() {
		conf.Comment = keygen.GetSSHKeyComment()
	} else {
		conf.Comment = data.Comment.String()
	}

	// data.Comment = types.StringValue("bar")

	if sshkey, err = keygen.New(&conf); err != nil {
		resp.Diagnostics.AddError("Key generation failed", err.Error())

		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.StringValue(sshkey.SHA256())
	data.PrivateKeyPEM = types.StringValue(string(sshkey.PrivateKeyPEM()))
	data.PublicKey = types.StringValue(string(sshkey.PublicKey()))
	data.FingerprintMD5 = types.StringValue(sshkey.MD5())
	data.FingerprintSHA256 = types.StringValue(sshkey.SHA256())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// no need to support Read at the moment since the resource is fully within state
// NOTE: if we support sourcing from files on disk in the future, this will have to be implemented.
func (r *SSHKeyPairResource) Read(
	_ context.Context,
	_ resource.ReadRequest,
	_ *resource.ReadResponse,
) {
}

// no need to support Read at the moment since the resource is fully within state
// NOTE: if we support sourcing from files on disk in the future, this will have to be implemented.
func (r *SSHKeyPairResource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	_ *resource.UpdateResponse,
) {
}

func (r *SSHKeyPairResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data *SSHKeyPairResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *SSHKeyPairResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
