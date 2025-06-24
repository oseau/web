import datetime

from sqlite_migrate import Migrations
from sqlite_utils.db import Database

# Pick a unique name here - it must not clash with other migration sets that
# the user might run against the same database.

migration = Migrations("count")


# Use this decorator against functions that implement migrations
@migration()
def create_table(db: Database):
    # db is a sqlite-utils Database instance
    db["count"].create(
        {"id": int, "count": int},
        pk="id",
    )


@migration()
def add_time(db: Database):
    # db is a sqlite-utils Database instance
    db["count"].add_column("time", datetime.datetime)


@migration()
def set_default_time(db: Database):
    # set the default time to the current timestamp
    db["count"].transform(defaults={"time": "CURRENT_TIMESTAMP"})


@migration()
def set_default_time_for_existing_rows(db: Database):
    # update all existing rows with timestamp set to 1970-01-01
    # this is not really a migration,
    # but a one-time operation to set the time for existing rows
    db.execute(
        "UPDATE count SET time = '1970-01-01 00:00:00' WHERE time IS NULL",
    )
