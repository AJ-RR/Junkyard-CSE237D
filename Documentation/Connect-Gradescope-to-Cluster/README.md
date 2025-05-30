# Connecting Gradescope to the cluster

The final step of the pipeline is connecting 

## Exposing the server externally

From the previous page ([GreenGrader JobServer Service](https://www.notion.so/GreenGrader-JobServer-Service-1f286fd69b8a8076ba68fcfca51c56ae?pvs=21)) we were able to access the job server from automaton, but 

**Requires `sudo` access on automaton**

Edit the apache2 config on automaton to create a reverse proxy for our job-server

```bash
sudo vim /etc/apache2/sites-available/smart-cycling.sysnet.ucsd.edu-le-ssl.conf
```

Inside the `<VirtualHost>` block, we want to include the following

```
ProxyPass /gradescope "http://[job-server ip]:30080" nocanon
ProxyPassReverse /gradescope "http://[job-server ip]:30080"
```

After saving the changes, restart the apache2 service with

```bash
sudo systemctl restart apache2
```

Now, the service should be accessible at [https://smartcycling.sysnet.ucsd.edu/gradescope](https://smartcycling.sysnet.ucsd.edu/gradescope)

![image.png](image.png)

## Building the Gradescope autograder

Now that the server can be accessed from outside the network, its possible to send requests from Gradescope

Building an autograder for Gradescope requires two bash script files

- setup-sh: to install all the required dependencies
- run_autograder: script that produces the testing results

**setup.sh**

```bash
#! /bin./bash
apt-get updateapt-get install -y curl jq zip
```

**run_autograder**

```bash
#!/bin/bashset -e
# Config
URL="https://smartcycling.sysnet.ucsd.edu/gradescope/submit"
SUBMISSION_DIR="/autograder/submission"
ZIP_FILE="/tmp/submission.zip"
RESULTS_JSON="/autograder/results/results.json"
METADATA_FILE="/autograder/submission_metadata.json"

STUDENT_NAME=$(jq -r '.users[0].name' "$METADATA_FILE" | tr ' ' '-')
ASSIGNMENT_TITLE=$(jq -r '.assignment.title' "$METADATA_FILE" | tr ' ' '-')

# Zip the submission files 
cd "$SUBMISSION_DIR"
zip -r "$ZIP_FILE" ./*cd -

# Make POST request, write logs from k8s job to /autograder/results/results.json
curl -X POST "$URL" \
  -F "name=$STUDENT_NAME" \
  -F "image=$ASSIGNMENT_TITLE" \
  -F "script=@$ZIP_FILE" \
  -H "Content-Type: multipart/form-data" \
  -o "$RESULTS_JSON"

# Output saved location
echo "Response written to $RESULTS_JSON"
```

**run_autograder** zips up the student’s submission files (in `/autograder/submission` ) and sends it in a POST request to the cluster. The response returned is written directly to `/autograder/results/results.json` , the JSON file that the web client reads from

**NOTE**: The job created by the server must **ONLY** write the results.JSON information to `stdout`, all other output must be suppressed

**NOTE:** The current implementation of the job server uses a ConfigMap to pass in the submission into the job. Because of the limitations of Kubernetes, this means that the zipped submission files **CANNOT** be ≥ 1MiB