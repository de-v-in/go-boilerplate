CREATE TABLE IF NOT EXISTS "articles" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "content" TEXT NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT timezone('UTC', NOW()) NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT timezone('UTC', NOW()) NOT NULL
);
