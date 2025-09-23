# Spring Boot Microservices Boilerplate

A production-ready microservices boilerplate built with Spring Boot 3, featuring JWT authentication, service discovery, API gateway, and database-per-service architecture.

## 🚀 Introduction

This boilerplate provides a complete foundation for building scalable microservices applications. It includes essential components like authentication, user management, content management, and service orchestration. Perfect for developers who want to start with a robust, well-structured microservices architecture without building everything from scratch.

## ✨ Features

### Core Services
- **🔐 Auth Service** - JWT authentication with refresh tokens
- **👤 User Service** - User profile management and CRUD operations
- **📝 Post Service** - Content management with CRUD operations
- **🌐 API Gateway** - Request routing, logging, and CORS handling
- **🔍 Discovery Server** - Service registry using Netflix Eureka

### Technology Stack
- **☕ Spring Boot 3.2.0** - Modern Java framework
- **🗄️ PostgreSQL** - Database-per-service architecture
- **🐳 Docker** - Containerization and orchestration
- **🔒 Spring Security** - Authentication and authorization
- **📊 Spring Data JPA** - Data persistence layer
- **🌊 Spring Cloud Gateway** - API gateway implementation

### Architecture Benefits
- **Database-per-service** pattern for true service isolation
- **JWT-based authentication** with refresh token support
- **Service discovery** for dynamic service registration
- **Request logging** and monitoring capabilities
- **CORS handling** for web application integration

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Applications                      │
└─────────────────────┬───────────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────────┐
│                    API Gateway (Port 8080)                     │
│  • Request Routing  • CORS Handling  • Request Logging        │
└─────────────────────┬───────────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
┌───────▼──────┐ ┌────▼────┐ ┌─────▼─────┐
│ Auth Service │ │Discovery│ │User Service│
│  (Port 8081) │ │(Port    │ │(Port 8082)│
│              │ │ 8761)   │ │           │
└───────┬──────┘ └─────────┘ └─────┬─────┘
        │                          │
        │                          │
┌───────▼──────┐              ┌────▼─────┐
│   authdb     │              │  userdb  │
│ PostgreSQL   │              │PostgreSQL│
└──────────────┘              └──────────┘
                                      │
                              ┌───────▼──────┐
                              │ Post Service │
                              │(Port 8083)   │
                              └───────┬──────┘
                                      │
                              ┌───────▼──────┐
                              │   postdb     │
                              │ PostgreSQL   │
                              └──────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Java 17+ (for local development)
- Maven 3.6+ (for local development)

### Running the Application

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd spring-boot-microservices
   ```

2. **Start all services with Docker Compose**
   ```bash
   docker-compose up --build
   ```

3. **Verify services are running**
   ```bash
   # Check service status
   docker-compose ps
   
   # Test API Gateway
   curl http://localhost:8080/api/auth/health
   ```

4. **Access the services**
   - **API Gateway**: http://localhost:8080
   - **Eureka Dashboard**: http://localhost:8761
   - **Auth Service**: http://localhost:8081
   - **User Service**: http://localhost:8082
   - **Post Service**: http://localhost:8083

### Development Mode

To run services locally for development:

```bash
# Build all services
mvn clean package -DskipTests

# Run Discovery Server
cd discovery-server && mvn spring-boot:run

# Run Auth Service (in another terminal)
cd auth-service && mvn spring-boot:run

# Run User Service (in another terminal)
cd user-service && mvn spring-boot:run

# Run Post Service (in another terminal)
cd post-service && mvn spring-boot:run

# Run API Gateway (in another terminal)
cd api-gateway && mvn spring-boot:run
```

## 📋 API Endpoints

### Authentication Service (`/api/auth`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register a new user |
| POST | `/api/auth/login` | User login |
| POST | `/api/auth/refresh` | Refresh JWT token |
| GET | `/api/auth/health` | Health check |

### User Service (`/api/users`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | Get all users |
| GET | `/api/users/{id}` | Get user by ID |
| POST | `/api/users` | Create new user |
| PUT | `/api/users/{id}` | Update user |
| DELETE | `/api/users/{id}` | Delete user |
| GET | `/api/users/health` | Health check |

### Post Service (`/api/posts`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/posts` | Get all posts |
| GET | `/api/posts/{id}` | Get post by ID |
| POST | `/api/posts` | Create new post |
| PUT | `/api/posts/{id}` | Update post |
| DELETE | `/api/posts/{id}` | Delete post |
| GET | `/api/posts/health` | Health check |

### Discovery Server
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/discovery/health` | Health check |
| GET | `/discovery/info` | Service information |

## 🔧 How to Extend

### Adding a New Service

1. **Create service directory structure**
   ```bash
   mkdir new-service
   cd new-service
   mkdir -p src/main/java/com/example/newservice
   mkdir -p src/main/resources
   ```

2. **Create `pom.xml`**
   ```xml
   <?xml version="1.0" encoding="UTF-8"?>
   <project xmlns="http://maven.apache.org/POM/4.0.0"
            xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
            xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
            http://maven.apache.org/xsd/maven-4.0.0.xsd">
       <modelVersion>4.0.0</modelVersion>
       
       <parent>
           <groupId>com.example</groupId>
           <artifactId>spring-boot-microservices</artifactId>
           <version>1.0.0</version>
       </parent>
       
       <artifactId>new-service</artifactId>
       <packaging>jar</packaging>
       
       <dependencies>
           <dependency>
               <groupId>org.springframework.boot</groupId>
               <artifactId>spring-boot-starter-web</artifactId>
           </dependency>
           <dependency>
               <groupId>org.springframework.cloud</groupId>
               <artifactId>spring-cloud-starter-netflix-eureka-client</artifactId>
           </dependency>
       </dependencies>
   </project>
   ```

3. **Create main application class**
   ```java
   package com.example.newservice;
   
   import org.springframework.boot.SpringApplication;
   import org.springframework.boot.autoconfigure.SpringBootApplication;
   
   @SpringBootApplication
   public class NewServiceApplication {
       public static void main(String[] args) {
           SpringApplication.run(NewServiceApplication.class, args);
       }
   }
   ```

4. **Create `application.yml`**
   ```yaml
   server:
     port: 8084
   
   spring:
     application:
       name: new-service
   
   eureka:
     client:
       service-url:
         defaultZone: http://localhost:8761/eureka/
     instance:
       prefer-ip-address: true
   ```

5. **Add service to parent `pom.xml`**
   ```xml
   <modules>
       <module>discovery-server</module>
       <module>api-gateway</module>
       <module>auth-service</module>
       <module>user-service</module>
       <module>post-service</module>
       <module>new-service</module>  <!-- Add this line -->
   </modules>
   ```

6. **Add service to `docker-compose.yml`**
   ```yaml
   new-service:
     build: ./new-service
     container_name: new-service
     ports:
       - "8084:8084"
     environment:
       - SPRING_PROFILES_ACTIVE=docker
       - EUREKA_CLIENT_SERVICE_URL_DEFAULTZONE=http://discovery-server:8761/eureka/
     networks:
       - microservices-network
     depends_on:
       discovery-server:
         condition: service_started
   ```

7. **Add route to API Gateway**
   ```yaml
   # In api-gateway/src/main/resources/application.yml
   spring:
     cloud:
       gateway:
         routes:
           - id: new-service
             uri: lb://new-service
             predicates:
               - Path=/api/new/**
             filters:
               - StripPrefix=1
   ```

### Adding a New Database

1. **Create database setup script**
   ```sql
   -- database-scripts/04-newdb-schema.sql
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
   
   CREATE TABLE IF NOT EXISTS new_entities (
       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
       name VARCHAR(100) NOT NULL,
       description TEXT,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

2. **Update `docker-compose.yml`**
   ```yaml
   postgres:
     volumes:
       - ./database-scripts/04-newdb-schema.sql:/docker-entrypoint-initdb.d/04-newdb-schema.sql
   ```

3. **Create database user**
   ```sql
   CREATE USER new_service_user WITH PASSWORD 'new_service_password';
   GRANT CONNECT ON DATABASE newdb TO new_service_user;
   GRANT USAGE ON SCHEMA public TO new_service_user;
   GRANT CREATE ON SCHEMA public TO new_service_user;
   ```

## 🛠️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SPRING_PROFILES_ACTIVE` | Spring profile | `docker` |
| `EUREKA_CLIENT_SERVICE_URL_DEFAULTZONE` | Eureka server URL | `http://discovery-server:8761/eureka/` |
| `SPRING_DATASOURCE_URL` | Database connection URL | Service-specific |
| `SPRING_DATASOURCE_USERNAME` | Database username | Service-specific |
| `SPRING_DATASOURCE_PASSWORD` | Database password | Service-specific |

### Database Configuration

Each service connects to its own PostgreSQL database:

- **Auth Service**: `authdb` (port 5432)
- **User Service**: `userdb` (port 5432)
- **Post Service**: `postdb` (port 5432)

## 📊 Monitoring

### Health Checks
- **Service Health**: `GET /actuator/health`
- **Service Info**: `GET /actuator/info`
- **Eureka Dashboard**: http://localhost:8761

### Logging
- **Request Logging**: All requests are logged through the API Gateway
- **Response Time**: Tracked and logged for performance monitoring
- **Service Registration**: Eureka logs all service registrations

## 🔒 Security

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: BCrypt password encryption
- **CORS Configuration**: Cross-origin resource sharing setup
- **Service Isolation**: Database-per-service for data isolation

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Credits

- **Spring Boot Team** - For the amazing framework
- **Netflix OSS** - For Eureka service discovery
- **PostgreSQL Team** - For the robust database
- **Docker Team** - For containerization platform

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

If you have any questions or need help with this boilerplate:

- 📧 Create an issue in this repository
- 📖 Check the documentation in each service directory
- 🔍 Review the example implementations

---

**Happy coding! 🚀**
