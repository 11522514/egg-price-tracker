# Egg Price Tracker

A comprehensive web application for tracking egg prices over time, comparing national averages with local prices, and visualizing historical data trends.

## Features

- **Historical Price Tracking**: Store and visualize egg price data over time
- **National vs Local Comparison**: Compare local prices against national averages
- **Interactive Charts**: Dynamic price visualization with Chart.js
- **Data Collection**: Automated data collection from multiple sources
- **REST API**: Full API for price data management
- **Real-time Updates**: Add new price data through the web interface

## Technology Stack

- **Backend**: Go (Golang) with Gorilla Mux
- **Database**: PostgreSQL
- **Data Collection**: Python with automated scheduling
- **Frontend**: HTML, CSS, JavaScript with Chart.js
- **Containerization**: Docker and Docker Compose

## Project Structure

```
egg-price-tracker/
├── main.go                 # Go backend server
├── go.mod                  # Go module dependencies
├── data_collector.py       # Python data collection service
├── scheduler.py            # Python scheduling service
├── requirements.txt        # Python dependencies
├── schema.sql              # Database schema
├── docker-compose.yml      # Docker orchestration
├── Dockerfile.go           # Go service Docker image
├── Dockerfile.python       # Python service Docker image
├── static/
│   └── index.html          # Web frontend
└── README.md               # This file
```

## Quick Start with Docker

1. **Clone and setup:**

   ```bash
   mkdir egg-price-tracker
   cd egg-price-tracker
   # Copy all the provided files to this directory
   ```

2. **Start the application:**

   ```bash
   docker-compose up -d
   ```

3. **Access the application:**
   - Web Interface: http://localhost:8080
   - API Base URL: http://localhost:8080/api

## Manual Setup (Development)

### Prerequisites

- Go 1.21+
- Python 3.11+
- PostgreSQL 15+
- Node.js (optional, for advanced frontend development)

### Database Setup

1. **Install PostgreSQL and create database:**

   ```bash
   # On Ubuntu/Debian
   sudo apt install postgresql postgresql-contrib

   # On macOS
   brew install postgresql

   # Create database
   sudo -u postgres createdb egg_tracker
   ```

2. **Run database schema:**
   ```bash
   sudo -u postgres psql egg_tracker < schema.sql
   ```

### Backend Setup (Go)

1. **Install dependencies:**

   ```bash
   go mod init egg-price-tracker
   go mod tidy
   ```

2. **Set environment variables:**

   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=your_password
   export DB_NAME=egg_tracker
   export PORT=8080
   ```

3. **Run the Go backend:**
   ```bash
   go run main.go
   ```

### Data Collection Setup (Python)

1. **Install Python dependencies:**

   ```bash
   pip install -r requirements.txt
   ```

2. **Set environment variables:**

   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=your_password
   export DB_NAME=egg_tracker
   ```

3. **Run data collection:**

   ```bash
   # One-time collection
   python data_collector.py

   # Scheduled collection
   python scheduler.py
   ```

## API Endpoints

### GET /api/prices

Get price data with optional filters.

**Query Parameters:**

- `location` (optional): Filter by location (e.g., "NATIONAL", "California")
- `limit` (optional): Number of records to return (default: 30)

**Example:**

```bash
curl "http://localhost:8080/api/prices?location=NATIONAL&limit=10"
```

### POST /api/prices

Add new price data.

**Request Body:**

```json
{
  "date": "2024-01-15",
  "location": "California",
  "price_per_dozen": 2.45,
  "source": "Local Market"
}
```

**Example:**

```bash
curl -X POST http://localhost:8080/api/prices \
  -H "Content-Type: application/json" \
  -d '{"date":"2024-01-15","location":"California","price_per_dozen":2.45,"source":"Local Market"}'
```

### GET /api/locations

Get all available locations.

### GET /api/comparison

Get price comparison between locations and national average.

## Configuration

### Environment Variables

| Variable      | Description       | Default       |
| ------------- | ----------------- | ------------- |
| `DB_HOST`     | Database host     | `localhost`   |
| `DB_PORT`     | Database port     | `5432`        |
| `DB_USER`     | Database username | `postgres`    |
| `DB_PASSWORD` | Database password | `password`    |
| `DB_NAME`     | Database name     | `egg_tracker` |
| `PORT`        | Server port       | `8080`        |

### Data Collection Configuration

The Python data collector can be configured to collect from various sources:

1. **USDA/FRED Data**: Set `FRED_API_KEY` environment variable
2. **Grocery Store APIs**: Configure store-specific API keys
3. **Collection Schedule**: Modify `scheduler.py` for different intervals

## Deployment Options

### Docker Deployment (Recommended)

1. **Production with Docker Compose:**

   ```bash
   # Use production configuration
   docker-compose -f docker-compose.prod.yml up -d
   ```

2. **Individual container deployment:**

   ```bash
   # Build images
   docker build -f Dockerfile.go -t egg-tracker-backend .
   docker build -f Dockerfile.python -t egg-tracker-collector .

   # Run containers
   docker run -d --name postgres -e POSTGRES_DB=egg_tracker postgres:15
   docker run -d --name backend --link postgres egg-tracker-backend
   docker run -d --name collector --link postgres egg-tracker-collector
   ```

### Cloud Deployment

#### AWS Deployment

1. **Use AWS RDS for PostgreSQL**
2. **Deploy Go backend to AWS ECS or Lambda**
3. **Run Python collector as scheduled ECS task or Lambda function**
4. **Use CloudFront for static file serving**

#### Google Cloud Platform

1. **Use Cloud SQL for PostgreSQL**
2. **Deploy to Cloud Run or GKE**
3. **Use Cloud Scheduler for data collection**

#### Digital Ocean

1. **Use managed PostgreSQL database**
2. **Deploy to App Platform or Droplets**
3. **Use cron jobs for data collection**

### Traditional Server Deployment

1. **Setup reverse proxy (Nginx):**

   ```nginx
   server {
       listen 80;
       server_name yourdomain.com;

       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

2. **Setup systemd services:**

   ```ini
   [Unit]
   Description=Egg Price Tracker Backend
   After=network.target

   [Service]
   Type=simple
   User=ubuntu
   WorkingDirectory=/home/ubuntu/egg-price-tracker
   ExecStart=/home/ubuntu/egg-price-tracker/egg-tracker
   Restart=always

   [Install]
   WantedBy=multi-user.target
   ```

## Monitoring and Maintenance

### Health Checks

- Backend health: `GET http://localhost:8080/api/prices?limit=1`
- Database connectivity: Check through application logs
- Data freshness: Monitor last collection timestamps

### Log Monitoring

- Go backend: Logs to stdout/stderr
- Python collector: Structured logging with timestamps
- Database: PostgreSQL logs for query performance

### Backup Strategy

```bash
# Daily database backup
pg_dump egg_tracker > backup_$(date +%Y%m%d).sql

# Automated backup script
#!/bin/bash
pg_dump egg_tracker | gzip > /backups/egg_tracker_$(date +%Y%m%d_%H%M%S).sql.gz
find /backups -name "*.sql.gz" -mtime +30 -delete
```

## Customization

### Adding New Data Sources

1. Extend `data_collector.py` with new collection methods
2. Add source-specific configuration
3. Update database schema if needed

### Frontend Customization

1. Modify `static/index.html` for UI changes
2. Add new chart types or visualizations
3. Implement additional filters or views

### API Extensions

1. Add new endpoints in `main.go`
2. Implement additional query parameters
3. Add data validation and error handling

## Troubleshooting

### Common Issues

1. **Database Connection Failed:**

   - Check PostgreSQL service status
   - Verify connection parameters
   - Ensure database exists

2. **Port Already in Use:**

   - Change PORT environment variable
   - Check for conflicting services

3. **Data Collection Errors:**
   - Check API keys and rate limits
   - Verify network connectivity
   - Review collector logs

### Debug Mode

```bash
# Enable debug logging
export LOG_LEVEL=debug

# Run with verbose output
go run main.go -v
python data_collector.py --debug
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes and test thoroughly
4. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
