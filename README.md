# AKS Orchestration

A web-based tool for orchestrating AKS (Azure Kubernetes Service) installation and configuration scripts. Built with Go and containerized for deployment in Azure Web App for Containers.

## Features

- Web interface for monitoring script execution
- Sequential script execution with status tracking
- Real-time execution status updates
- Error handling and reporting
- Azure DevOps pipeline integration
- Containerized deployment

## Project Structure

```
.
├── main.go                 # Main application file
├── templates/
│   └── index.html         # Web interface template
├── scripts/
│   ├── 1_aks_prerequisites.sh
│   ├── 2_create_aks_cluster.sh
│   ├── 3_setup_networking.sh
│   └── 4_install_monitoring.sh
├── Dockerfile             # Container configuration
├── azure-pipelines.yml    # Azure DevOps pipeline configuration
└── go.mod                 # Go module file
```

## Local Development

1. Prerequisites:
   - Go 1.21 or later
   - Docker (optional)

2. Build and Run:
   ```bash
   # Build
   go mod download
   go build -o script-orchestrator

   # Run
   ./script-orchestrator
   ```

3. Access the web interface at `http://localhost:8080`

## Docker Build

```bash
docker build -t script-orchestrator .
docker run -p 8080:8080 script-orchestrator
```

## Azure DevOps Deployment

1. Configure Azure DevOps pipeline variables:
   - `AZURE_SUBSCRIPTION`
   - `RESOURCE_GROUP`
   - `APP_SERVICE_NAME`

2. Set up service connections:
   - Azure Resource Manager connection
   - Container Registry connection

3. Create Azure resources:
   ```bash
   az group create --name your-resource-group --location your-location
   az appservice plan create --name script-orchestrator-plan --resource-group your-resource-group --sku B1 --is-linux
   az webapp create --resource-group your-resource-group --plan script-orchestrator-plan --name your-web-app-name --deployment-container-image-name your-container-registry/script-orchestrator:latest
   ```

4. Run the pipeline in Azure DevOps

## Script Customization

The scripts in the `scripts/` directory can be customized for your specific AKS installation requirements. Each script should:
- Be executable
- Return appropriate exit codes
- Handle errors appropriately
- Provide meaningful output

## License

MIT License 