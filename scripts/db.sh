#!/bin/bash

# Database management script for Wordle

set -e

COMPOSE_FILE="docker-compose.yml"
DB_SERVICE="postgres"

case "$1" in
    "start")
        echo "Starting PostgreSQL database..."
        docker-compose up -d $DB_SERVICE
        echo "Database started. Waiting for it to be ready..."
        sleep 5
        docker-compose exec $DB_SERVICE pg_isready -U wordle_user -d wordle
        echo "Database is ready!"
        ;;
    "stop")
        echo "Stopping PostgreSQL database..."
        docker-compose stop $DB_SERVICE
        ;;
    "restart")
        echo "Restarting PostgreSQL database..."
        docker-compose restart $DB_SERVICE
        ;;
    "reset")
        echo "Resetting PostgreSQL database (this will delete all data)..."
        read -p "Are you sure? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker-compose down -v
            docker-compose up -d $DB_SERVICE
            echo "Database reset complete!"
        else
            echo "Reset cancelled."
        fi
        ;;
    "shell")
        echo "Connecting to PostgreSQL shell..."
        docker-compose exec $DB_SERVICE psql -U wordle_user -d wordle
        ;;
    "logs")
        echo "Showing PostgreSQL logs..."
        docker-compose logs -f $DB_SERVICE
        ;;
    "status")
        echo "Checking database status..."
        docker-compose ps $DB_SERVICE
        ;;
    "backup")
        BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
        echo "Creating backup: $BACKUP_FILE"
        docker-compose exec $DB_SERVICE pg_dump -U wordle_user wordle > "$BACKUP_FILE"
        echo "Backup created: $BACKUP_FILE"
        ;;
    "restore")
        if [ -z "$2" ]; then
            echo "Usage: $0 restore <backup_file>"
            exit 1
        fi
        echo "Restoring from backup: $2"
        docker-compose exec -T $DB_SERVICE psql -U wordle_user wordle < "$2"
        echo "Restore complete!"
        ;;
    *)
        echo "Wordle Database Management Script"
        echo ""
        echo "Usage: $0 {start|stop|restart|reset|shell|logs|status|backup|restore}"
        echo ""
        echo "Commands:"
        echo "  start    - Start the PostgreSQL database"
        echo "  stop     - Stop the PostgreSQL database"
        echo "  restart  - Restart the PostgreSQL database"
        echo "  reset    - Reset database (deletes all data)"
        echo "  shell    - Connect to PostgreSQL shell"
        echo "  logs     - Show database logs"
        echo "  status   - Show database status"
        echo "  backup   - Create a database backup"
        echo "  restore  - Restore from a backup file"
        echo ""
        exit 1
        ;;
esac
