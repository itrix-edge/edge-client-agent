Edge Client Agent
=================
An agent to deploy and manage applications for the kubernetes edge cluster.

## Use scenarios
1. Register application deployments to the agent
2. Hit trigger to start deployments and application
3. Hit tirgger again to perform deployment cycle

### Register application deployments
a. Users prepare application deployment and service YAML.
b. Users register the deployment and serivce YAML to the agent.
c. The agent reply registeration with the hook key, save the key for later use.

### Hit trigger to start deployments
a. Users design their continus integration flow, hit the hook URL with the saved key to the agent when needed.
b. Agent executes hook related application deployment and its service with registered information.
c. Local kubernetes cluster creates namespace, deployment and the service, start user applications.

### Hit tirgger again for deployment cycle
a. User's continus integration flow hit the hook again with different POST DATA included.
b. Agent perform application new version upgrade with given new POST DATA for the deployment.
c. Local kubernetes cluster uses rolling-upgrade with the deployment.

### How to use

1. Deploy db.
2. Deploy with YAML:
   ```=shell
   kubectl apply -f edge-agent-all.yaml
   ```
3. Enable migration.
4. Register Deployment and Services.
5. Hit agent with assoicated key to start deployment.