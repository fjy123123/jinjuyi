services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile.railway
    port: 8080
    healthCheck:
      path: /health
      interval: 30
      timeout: 10
      retries: 3
    environment:
      - MYSQL_HOST=${{MYSQL_HOST}}
      - MYSQL_PORT=${{MYSQL_PORT}}
      - MYSQL_USER=${{MYSQL_USER}}
      - MYSQL_PASSWORD=${{MYSQL_PASSWORD}}
      - MYSQL_DATABASE=${{MYSQL_DATABASE}}
      - MONGO_HOST=${{MONGO_HOST}}
      - MONGO_PORT=${{MONGO_PORT}}
      - MONGO_DB=${{MONGO_DB}}
      - REDIS_HOST=${{REDIS_HOST}}
      - REDIS_PORT=${{REDIS_PORT}}
      - REDIS_PASSWORD=${{REDIS_PASSWORD}}
      - JWT_SECRET=${{JWT_SECRET}}
      - GIN_MODE=release
    depends_on:
      - mysql
      - mongodb
      - redis

  mysql:
    image: mysql:8.0
    port: 3306
    environment:
      - MYSQL_ROOT_PASSWORD=${{MYSQL_ROOT_PASSWORD}}
      - MYSQL_DATABASE=${{MYSQL_DATABASE}}
      - MYSQL_USER=${{MYSQL_USER}}
      - MYSQL_PASSWORD=${{MYSQL_PASSWORD}}
    volumes:
      - mysql_data:/var/lib/mysql

  mongodb:
    image: mongo:7
    port: 27017
    volumes:
      - mongo_data:/data/db

  redis:
    image: redis:7-alpine
    port: 6379
    environment:
      - REDIS_PASSWORD=${{REDIS_PASSWORD}}
    volumes:
      - redis_data:/data

volumes:
  mysql_data:
  mongo_data:
  redis_data: