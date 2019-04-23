# Git Orchestrator - An operator for Kubernetes

## Why?
This operator allows you to provide the creation of repositories in an automated manner, on your terms.
You control the base-rules, others can create repositories.

## Setup

### Docker file
case 1: Build your own image.

`operator-sdk build whatever-repo` and change the image in `operator.yaml`

case 2: Use my image, which is already defined. 

### Settings

```
---
apiVersion: gitorchestrator.sysrant.com/v1alpha1
kind: GitOrchestrator
metadata:
  name: gitorchestrator
spec:
  gitType: "github"
  organisation: "leave-empty-if-none"
  repositoryName: "example-repo-k8s"
  repositoryNamespace: "wiardvanrij"
  description: "Created from k8s!"
  visibility: "public"
```

* gitType: either `github` or `gitlab` for now
* organisation: only for github, leave empty if it is a personal project
* repositoryName: the name of the repository
* repositoryNamespace: this is the "path" to your repository. Both github as gitlab use this. For example: /wiardvanrij/myrepo  - wiardvanrij is the namespace. This is also required for organisations.
* description: the description of the repository
* visibility: either `private` or `public` for github, gitlab also has `internal`

The `operator.yaml` needs a few more inputs:

* SECRET_TOKEN: which comes from the `secret.yaml` file. For `github` create an oauth token, for `gitlab` use a deployment token. Base64 this in the `secret.yaml`
* END_POINT: ignore this for `github` - Gitlab uses `https://gitlab.com/api/v4` or `v3` if you are oldscool. If you run gitlab yourself, find your own link ;)

### Deployment
Deploy the files in the deploy folder after you are done with the setup.

## More info

Even though the operator knows what has been processed, it does not "update" repositories. Ea; you cannot change its name or "visibility". Changing the name will never happen, visibility perhaps in the future. We also do nothing with deleted resources: so no, it does not delete your repository.  

## Known issues
* It can only process one resources at a time. - this has something to do with persisting the current status of the resources. c.b.a
* It's not fool proof 

## In production
I assume that if you are going to make use of this, that you have the knowledge of what you are doing and what the state of this operator is. The core of this project works fine, yet it could be more mature.

## Help out
So only github & gitlab are implemented. Feel free to make a PR to add others. 

Also feel free to make PR's to make this operator more mature.
