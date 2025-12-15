#!/bin/sh
set -e

echo "========================================="
echo "🚀 E-Commerce Backend Starting..."
echo "========================================="

echo ""
echo "📡 Checking PostgreSQL connection..."

# Wait for PostgreSQL
MAX_RETRIES=30
RETRY_COUNT=0

until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q' 2>/dev/null; do
  RETRY_COUNT=$((RETRY_COUNT + 1))
  
  if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
    echo "❌ Failed to connect to PostgreSQL after $MAX_RETRIES attempts"
    exit 1
  fi
  
  echo "⏳ PostgreSQL is unavailable - waiting... (attempt $RETRY_COUNT/$MAX_RETRIES)"
  sleep 2
done

echo "✅ PostgreSQL is ready!"
echo ""
echo "========================================="
echo "🚀 Starting Application..."
echo "========================================="

# Start the application
exec /app/main