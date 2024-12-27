# Use a lightweight base image for PostgreSQL
FROM postgres:alpine

# Set environment variables for PostgreSQL
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=123
ENV POSTGRES_DB=PackageTracking_db

RUN chmod -R 777 /var/lib/postgresql/data

# Expose the PostgreSQL port
EXPOSE 5432
