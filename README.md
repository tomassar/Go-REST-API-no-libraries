This is a small project creating a REST API using pure Go, using only the standard HTTP package.
Consists on a simple domain object called OpenSourceProject that tracks open source projects

Currently the server only supports these methods:

- Get all projects
- Post a project
- Get a project by id
- Get admin dashboard only if basic auth success

And implements the next exceptions:

- Return an error if the Content Type is not Application/JSON
- If trying to get admin dashboard and basic auth failed, then return unauthorized
