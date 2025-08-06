<h1 align="center">E-commerce Customers Service</h1>

<p align="center">
	<img src="docs/images/golang.jpg" width="150" height="150"/>
	<img src="docs/images/postgres.jpg" width="150" height="150"/>
	<img src="docs/images/redis.png" width="170" height="150"/>
	<img src="docs/images/prometheus.png" width="170" height="150">
	<img src="docs/images/grafana.png" width="170" height="150">
</p>

## Step 2: Install Dependencies

Navigate to the project's root directory and run the following command to install Go dependencies:

```bash
go mod download
```

## Step 3: Configure Environment Variables

The project may require some environment variables to work correctly. Check if there is a `.env.example` file in the root directory. If it exists, make a copy and rename it to `.env`. Then update the variables according to your local environment settings.

## Step 4: Start the Environment

To start the services, run the following command:

```bash
make up
```

## Step 4.1: Stop the Environment

To stop the services, run:

```bash
make down
```

## Step 5: Run the Application

To run the application, use the following command:  
**Note:** Go version **1.22.1 or higher** must be installed on your machine.

```bash
make run-app
```

or

```bash
make run-app-logs
```

## Step 7: Run Tests

To run unit tests, use:

```bash
make unit-tests
```

**Note:** The application must be running to execute integration tests.

```bash
make integration-tests
```

## Step 8: Generate a New Version Tag

To generate a new version tag, run the following command and follow the terminal instructions:

```bash
make tag
```

## Reference Articles

The following articles were used as references in the design and architecture of this project:

- [Explicit Architecture: DDD, Hexagonal, Onion, Clean, CQRS - How I Put It All Together](https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/)
- [Sharding & IDs at Instagram](https://instagram-engineering.com/sharding-ids-at-instagram-1cf5a71e5a5c)
