CREATE TABLE "Users" (
  "id" integer PRIMARY KEY,
  "username" varchar NOT NULL UNIQUE,
  "password" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "join_date" date
);

CREATE TABLE "UserProfiles" (
  "id" integer PRIMARY KEY,
  "user_id" integer NOT NULL,
  "avatar" varchar NOT NULL,
  "bio" varchar NOT NULL,
  "location" varchar NOT NULL DEFAULT 'Kazakhstan'
);

CREATE TABLE "Forums" (
  "id" integer PRIMARY KEY,
  "name" varchar NOT NULL UNIQUE,
  "description" varchar NOT NULL
);

CREATE TABLE "Topics" (
  "id" integer PRIMARY KEY,
  "forum_id" integer,
  "title" varchar NOT NULL UNIQUE,
  "created_by" integer
);

CREATE TABLE "Posts" (
  "id" integer PRIMARY KEY,
  "topic_id" integer NOT NULL,
  "user_id" integer NOT NULL,
  "body" text,
  "created_at" date
);

CREATE TABLE "Categories" (
  "id" integer PRIMARY KEY,
  "name" varchar NOT NULL UNIQUE
);

CREATE TABLE "ForumCategories" (
  "forum_id" integer NOT NULL,
  "category_id" integer NOT NULL
);

CREATE TABLE "TopicCategories" (
  "topic_id" integer NOT NULL,
  "category_id" integer NOT NULL
);

CREATE TABLE "Votes" (
  "id" integer PRIMARY KEY, 
  "user_id" integer NOT NULL,
  "post_id" integer NOT NULL,
  "vote_type" varchar NOT NULL
);

CREATE TABLE "Attachments" (
  "id" integer PRIMARY KEY,
  "post_id" integer NOT NULL,
	"file" text NOT NULL,
  "filename" varchar NOT NULL DEFAULT 'File',
  "size" integer
);

CREATE TABLE "Banlist" (
  "id" integer PRIMARY KEY,
  "user_id" integer NOT NULL,
  "banned_by" integer,
  "reason" varchar,
  "expires_at" date
);

ALTER TABLE "UserProfiles" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "Topics" ADD FOREIGN KEY ("forum_id") REFERENCES "Forums" ("id");

ALTER TABLE "Topics" ADD FOREIGN KEY ("created_by") REFERENCES "Users" ("id");

ALTER TABLE "Posts" ADD FOREIGN KEY ("topic_id") REFERENCES "Topics" ("id");

ALTER TABLE "Posts" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "ForumCategories" ADD FOREIGN KEY ("forum_id") REFERENCES "Forums" ("id");

ALTER TABLE "ForumCategories" ADD FOREIGN KEY ("category_id") REFERENCES "Categories" ("id");

ALTER TABLE "TopicCategories" ADD FOREIGN KEY ("topic_id") REFERENCES "Topics" ("id");

ALTER TABLE "TopicCategories" ADD FOREIGN KEY ("category_id") REFERENCES "Categories" ("id");

ALTER TABLE "Votes" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "Votes" ADD FOREIGN KEY ("post_id") REFERENCES "Posts" ("id");

ALTER TABLE "Attachments" ADD FOREIGN KEY ("post_id") REFERENCES "Posts" ("id");

ALTER TABLE "Banlist" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "Banlist" ADD FOREIGN KEY ("banned_by") REFERENCES "Users" ("id");