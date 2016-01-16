BEGIN;

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

CREATE TABLE "user" (
  "id" serial PRIMARY KEY,
  "username" citext NOT NULL,
  "password" text NOT NULL,
  "password_version" int NOT NULL DEFAULT 1
);

CREATE UNIQUE INDEX "user_username_idx" ON "user" ("username");

CREATE TABLE "log" (
  "id" serial PRIMARY KEY,
  "user_id" int NOT NULL,
  "name" text NOT NULL,
  "gpx" xml NOT NULL,
  "start" timestamptz NOT NULL,
  "end" timestamptz NOT NULL,
  "duration" int NOT NULL,
  "distance" int NOT NULL,
  "created" timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX "log_name_idx" ON "log" ("name");
CREATE INDEX "log_start_idx" ON "log" ("start");
CREATE INDEX "log_end_idx" ON "log" ("end");
CREATE INDEX "log_duration_idx" ON "log" ("duration");
CREATE INDEX "log_distance_idx" ON "log" ("distance");

CREATE TABLE "track" (
  "id" serial PRIMARY KEY,
  "log_id" int NOT NULL,
  "name" text,
  "start" timestamptz NOT NULL,
  "end" timestamptz NOT NULL,
  "duration" int NOT NULL,
  "distance" int NOT NULL,
  FOREIGN KEY ("log_id") REFERENCES "log" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX "track_start_idx" ON "track" ("start");
CREATE INDEX "track_end_idx" ON "track" ("end");
CREATE INDEX "track_duration_idx" ON "track" ("duration");
CREATE INDEX "track_distance_idx" ON "track" ("distance");

CREATE TABLE "trackpoint" (
  "id" serial PRIMARY KEY,
  "track_id" int NOT NULL,
  "point" point NOT NULL,
  "time" timestamptz NOT NULL,
  "elevation" real,
  "heartrate" int,
  "deleted" boolean NOT NULL DEFAULT false,
  FOREIGN KEY ("track_id") REFERENCES "track" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX "trackpoint_time_idx" ON "trackpoint" ("time");
CREATE INDEX "trackpoint_deleted_idx" ON "trackpoint" ("deleted");

CREATE TABLE "log_tag" (
  "log_id" int NOT NULL,
  "tag" citext NOT NULL,
  PRIMARY KEY ("log_id", "tag"),
  FOREIGN KEY ("log_id") REFERENCES "log" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX "log_tag_tag_idx" ON "log_tag" ("tag");

COMMIT;
