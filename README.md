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

1. Deploy database instance. Accroding to the deploy environment, its recommand that use external database service. For resource limited edge clusters, we provide following two alternative ways:
   a. For minimal performance impact the database can be run as external container, and use external services in kubernetes cluster for agent access:
		```=shell
		// Run containerized postgresql database instance. Note: data will lose after container terminated!
		$ docker run -it --name postgresql-local -p [Access_IP]:5432:5432/tcp -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgresdb -d postgres:11.8
		```
		* Please note that Access_IP should be assible IP address on host machine.
   b. Modify `external-db.service.yaml`, change `[Access_IP]` to your fit your deploy environment.
   c. It can also deploy postgresql services onto local Kubernetes cluster.
2. Deploy with YAML:
   ```=shell
   $ kubectl apply -f edge-agent-all.yaml
   ```
3. Enable migration. Touch server to perform migration.
   ```=shell
   $ curl http://[cluster-IP]:9000/v1/migrate/hook
   $ curl http://[cluster-IP]:9000/v1/migrate/deploymentTemplate
   ```
4. Register Deployment and Services.
   ```=shell
   $ curl -X POST http://[cluster-IP]:9000/v1/deploymentTemplate -d '{ "namespace": "Customizable", "options": [{"key": "template_key", "value": "template_value" }], "deployment_template": {...}, "service_template": {...}  }'
   ```
5. Hit agent with assoicated key to start deployment.