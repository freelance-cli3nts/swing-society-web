# Variables
IMAGE_NAME="europe-north1-docker.pkg.dev/swingsociety-backend/ss-go/ss-go"
SERVICE_NAME="ss-go"
REGION="europe-north1"  # Changed to eu region as per your request
ALLOW_UNAUTHENTICATED=false  # Set to false to stop unauthenticated access



# Step 1: Build the Docker container
echo "Building Docker container..."
docker build -t $IMAGE_NAME .

if [ $? -ne 0 ]; then
  echo "Docker build failed. Exiting."
  exit 1
fi

# Step 2: Push the Docker image to Google Container Registry
echo "Pushing Docker image to Google Container Registry..."
docker push $IMAGE_NAME

if [ $? -ne 0 ]; then
  echo "Docker push failed. Exiting."
  exit 1
fi

# Step 3: Deploy to Google Cloud Run
echo "Deploying to Google Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --project=swingsociety-backend \
  --service-account=ss-go1@swingsociety-backend.iam.gserviceaccount.com \
  --no-allow-unauthenticated  # Disable unauthenticated access

if [ $? -ne 0 ]; then
  echo "Cloud Run deployment failed. Exiting."
  exit 1
fi

echo "Deployment completed successfully!"