#aws-ssh

Shell into any server in your aws account.

#Search
calling the search command will look up all ec2 instances then filter the results by the term entered. if only a single result the app will ssh directly. If multiple results returned the app will display a numbered list and prompt for which to connect.

```
aws-ssh search {{search term}}
```

#Connect
connect allows you to connect directly to a known server. (team member running windows asked for this.)

#Configuration

you  will need to specify aws credentials. see [AWS Credentials](http://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs) for details.


In the central credentials file (~/.aws/credentials or %USERPROFILE%\.aws\credentials).
```
[default]
aws_access_key_id = ACCESS_KEY
aws_secret_access_key = SECRET_KEY
```

##PEM keys
if setting a key add the name only expected to be at ~/.ssh/{{KEY_NAME}}.pem
