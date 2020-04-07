# delphisbe
Backend monolith for delphis

Note: To get service running locally, kick off enginx using the nginx/.../local.delphishq.com config. Also need to add mapping to /etc/hosts for local.delphishq.com

# Get started
A few things you'll need to install. Assuming you're on mac:
* Make sure you have homebrew installed. If not run: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"`
* Install golang: `brew install go`
* Install nginx: `brew install nginx`
* Install local dynamodb `https://s3-us-west-2.amazonaws.com/dynamodb-local/dynamodb_local_latest.tar.gz` and unzip it to some place you can find it (e.g. ~/dynamodb/).
  * Run dynamodb: `java -Djava.library.path=<wherever you unzipped it> -jar DynamoDBLocal.jar -sharedDb`
  * Create the local tables `go run tools/create_dynamo/create_local_dynamo_tables.go`
* Setup internal dependencies: `make setup-internal-dep`
* Add an entry into your `/etc/hosts` that looks like `127.0.0.1	local.delphishq.com`. This is required for cookie management (cookies don't like localhost).
* Get the twitter credentials and put them in your environment:
    * `export DELPHIS_TWITTER_CONSUMER_KEY=xxxxx`
    * `export DELPHIS_TWITTER_CONSUMER_SECRET=xxxxx`
* Create an hmac secret (can be anything) and put it in your environment
    * `export DELPHIS_AUTH_HMAC_SECRET=xxxxx`
* Start nginx locally using the config: `make local-nginx`
* Start the go server: `go run server.go`

Okay at this point you should have the local server up and running. It has no data so let's do that:
* Login by going to `http://local.delphishq.com:8000/twitter/login`. This will actually log you in via twitter.
* Now, create a new discussion by going to `http://local.delphishq.com:8000/graphiql` and enter the following code:
```
mutation createDiscussion{
  createDiscussion(anonymityType:WEAK){
    id
  }
}
```
* Verify the discussionw as created by listing all discussions:
```
query listDiscussions{
  listDiscussions{
    id
    moderator{
      userProfile{
        twitterURL{
          url
        }
      }
    }
  }
}
```
This should show your name as the moderator.
* Now, start the client by navigating to `cd ./delphis-demo` and run `yarn start`. This should start the local app and allow you to navigate to `http://local.delphishq.com:8000/`.

And you're done!

# Build and deploy docker image
Retrieve credentials, build image, tag image, and push image.

```
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe
docker build -t delphisbe .
docker tag delphisbe:latest 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe:latest
docker push 033236388136.dkr.ecr.us-west-2.amazonaws.com/delphisbe:latest
```

# Database setup
I had to do the following:
* Create the database cluster within the console in order to not store the password (left a note in the aurora.tf file in terraform.)
* Create an instance (created via terraform)
* Create a bastion box ssh through (alternatively could add my IP to the VPC)
* The bastion box needs to be added to the RDS security group.
* Once on the db instance:
> CREATE USER chatham_staging WITH PASSWORD <REDACTED>;
> CREATE ROLE chatham_app;
> GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public to chatham_app;
> GRANT chatham_app TO chatham_staging;