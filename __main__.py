import pathlib
import pulumi
import pulumiverse_scaleway as scaleway
from pulumiverse_scaleway.iam import PolicyRuleArgs
import pulumi_pulumiservice as pulumiservice

from pydantic import BaseModel

class ProjectInfo(BaseModel):
    name: str
    description: str

config = pulumi.Config()

# From ESC environment
scw_organization_id = config.require_secret("organizationId")

# From stack local config
project_info = ProjectInfo(**config.require_object("project"))

# Scaleway project
project = scaleway.account.Project(
    "scaleway-project",
    name=project_info.name,
    description=project_info.description,
    organization_id=scw_organization_id,
)

# Scaleway Admin Group
admin_group = scaleway.iam.Group(
    "admin-group",
    name=f"{project_info.name}-administrators",
    organization_id=scw_organization_id,
)

policy = scaleway.iam.Policy(
    "policy-admin",
    name=f"{project_info.name}-admin",
    group_id=admin_group.id,
    organization_id=scw_organization_id,
    rules=[
        PolicyRuleArgs(
            organization_id=scw_organization_id,
            permission_set_names=["IAMManager"]
        ),
        PolicyRuleArgs(
            project_ids=[project.id],
            permission_set_names=["AllProductsFullAccess"]
        ),
    ],
)

app = scaleway.iam.Application(
    "app-pulumi",
    name=f"{project_info.name}-pulumi",
    organization_id=scw_organization_id,
)

membership = scaleway.iam.GroupMembership(
    "admin-group-member-app-pulumi",
    group_id=admin_group.id,
    application_id=app.id,
)

api_key = scaleway.iam.ApiKey(
    "apikey-pulumi",
    description=f"Pulumi API Key for {project_info.name}-pulumi",
    application_id=app.id,
    default_project_id=project.id,
)

# Create an ESC environment fetching the stack outputs of this same stack.
esc_template_path = pathlib.Path(__file__).parent / "environment.yaml"
esc_template = esc_template_path.read_text()
esc_actual_env = esc_template.format(stack=pulumi.get_stack())

escEnvironment = pulumiservice.Environment(
    "environmentResource",
    name=project.name,
    organization=pulumi.get_organization(),
    yaml=pulumi.StringAsset(esc_actual_env),
    project="scaleway"
)

pulumi.export("organizationId", scw_organization_id)
pulumi.export("projectId", project.id)
pulumi.export("accessKey", api_key.access_key)
pulumi.export("secretKey", api_key.secret_key)
with open('./Pulumi.README.md') as f:
    pulumi.export('readme', f.read())
