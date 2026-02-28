# Scaleway Projects

This Pulumi project sets up separate Scaleway projects. Each Scaleway project is a separate Pulumi stack
and contains the basic project setup

* Scaleway project ${project.name} (id: ${outputs.projectId})
* IAM
    * Administrators group ${admin_group.name}
    * Policy ${policy.name} allowing for:
      * Full IAM access to the organization (it's not possible to restrict access to a specific project)
      * Full access to all Scaleway services within the project
    * Application ${app.name} 