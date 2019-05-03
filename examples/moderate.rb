database :postgres do
  host '127.0.0.1'
  port 5432
  name :my_database
  #schema_search :public
end

server do
  host :localhost
  port 2727
end

jobs do
  server_uri 'http://127.0.0.1:9012'
  concurrency 10
end
