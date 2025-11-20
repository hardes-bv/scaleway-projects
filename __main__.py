import pulumi
import pulumiverse_scaleway as scaleway
from pulumiverse_scaleway.iam import PolicyRuleArgs
import pulumi_pulumiservice as puservice

config = pulumi.Config()

project_name = "hardes"
organization_id = config.require_secret("organizationId")

project = scaleway.account.Project(
    project_name,
    name=project_name,
    description="Hardes",
    organization_id=organization_id,
)

admin_group = scaleway.iam.Group(
    f"{project_name}-dns",
    organization_id=organization_id,
)

policy = scaleway.iam.Policy(
    f"{project_name}-admin",
    name=f"{project_name}-admin",
    group_id=admin_group.id,
    organization_id=organization_id,
    rules=[
        PolicyRuleArgs(
            organization_id=organization_id,
            permission_set_names=["IAMManager"]
        ),
        PolicyRuleArgs(
            project_ids=[project.id],
            permission_set_names=["AllProductsFullAccess"]
        ),
    ])

app = scaleway.iam.Application(
    "pulumi-hardes",
    organization_id=organization_id,
)

membership = scaleway.iam.GroupMembership(
    "pulumi-hardes-member",
    group_id=admin_group.id,
    application_id=app.id,
)

api_key = scaleway.iam.ApiKey(
    "pulumi-hardes",
    application_id=app.id,
)

hardesEscEnvironment = puservice.Environment(
    "environmentResource",
    name="hardes",
    organization=pulumi.get_organization(),
    yaml=pulumi.FileAsset("environment.yaml"),
    project="scaleway"
)

pulumi.export("organizationId", organization_id)
pulumi.export("projectId", project.id)
pulumi.export("accessKey", api_key.access_key)
pulumi.export("secretKey", api_key.secret_key)
