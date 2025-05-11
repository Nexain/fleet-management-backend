# Fleet Management Backend

This project is a backend service for managing vehicle fleets, designed to receive vehicle location data, store it in a PostgreSQL database, and provide APIs for accessing the latest location and travel history of vehicles. It also utilizes RabbitMQ for event handling related to geofencing.

## Table of Contents
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Docker Setup](#docker-setup)
- [Testing](#testing)
- [Contributing](#contributing)

## Technologies Used
- **Golang**: The primary programming language for backend development.
- **Eclipse Mosquitto**: MQTT broker for receiving vehicle location data.
- **PostgreSQL**: Database for storing vehicle location data.
- **RabbitMQ**: Message broker for handling geofence events.
- **Docker**: Containerization tool for running the application and its dependencies.

## Getting Started
To run the application locally, follow these steps:

1. **Clone the repository**:
   ```
   git clone https://github.com/yourusername/fleet-management-backend.git
   cd fleet-management-backend
   ```

2. **Build the Docker image**:
   ```
   docker build -t fleet-management-backend .
   ```

3. **Run the application using Docker Compose**:
   ```
   docker-compose up
   ```

4. **Access the API**:
   The API will be available at `http://localhost:8080`.

## API Endpoints
- **Get Last Location**:
  - **Endpoint**: `GET /vehicles/{vehicle_id}/location`
  - **Description**: Retrieves the last known location of the specified vehicle.

- **Get Location History**:
  - **Endpoint**: `GET /vehicles/{vehicle_id}/history?start={start_timestamp}&end={end_timestamp}`
  - **Description**: Retrieves the location history of the specified vehicle within the given time range.

## Docker Setup
The project includes a `docker-compose.yaml` file that sets up the following services:
- Backend service
- PostgreSQL database
- RabbitMQ message broker
- MQTT broker (Eclipse Mosquitto)

To start all services, run:
```
docker-compose up
```

## Testing
You can use Postman to test the API endpoints. A Postman collection is included in the repository for your convenience.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.