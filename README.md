# Mentat Analytics backend
Restful Backend Repo of Mentat Analytics.

### Microservices
* Dune Microservices: the microservice that handles the Rest API request with the Dune Analytics (requests the response by given auto-generated SQL query)
* OpenAI Microservices: the microservice that handles the Rest API request with the OpenAI Analytics (requests the auto-generated SQL query by given input)
* API Gateway: token generation/deciphering features, manages AWS SQS (message queue), querying the user requests history, burdening a limit to requests (up to 4)...
  ![image](https://user-images.githubusercontent.com/41055141/211517716-5b2873dc-1fe3-4324-b4ef-ac3d1316869b.png)

### What's this?
* A telegram chatbot that automates the process of doing the following:
    * creates SQL query using ChatGPT from the given database schema
      ![image](https://user-images.githubusercontent.com/41055141/211511821-b4271306-c826-4fc5-85c4-eee004da3f85.png)
      ![image](https://user-images.githubusercontent.com/41055141/211511850-beb55f11-3011-4228-947e-d761f7f6bc8b.png)
    * then simultaneously ask the query to Dune Analytics
      ![image](https://user-images.githubusercontent.com/41055141/211512766-c97f183b-242a-4c35-98d5-d66f91df3921.png)

### Tech Stacks
* Language: Go 1.19
* Web Framework: [Echo](https://echo.labstack.com/guide/customization/)
* DB Handler: [gorm](https://gorm.io/)
* Database: AWS RDS PostgreSQL
* Cache: Redis
* Logging: [Uber Zap](https://github.com/uber-go/zap)
* Database Migrations: [golang-migrate](https://github.com/golang-migrate/migrate)
* Unit/Mock & E2E Testing: go test, [testify](https://github.com/stretchr/testify), [newman](https://github.com/postmanlabs/newman)
* Configuration Management: [konaf](https://github.com/knadh/koanf)
* Deployment: Kubernetes (Jenkins/AWS kops, AWS CodeBuild)

### User Scenario & UI Mock
* [See here](https://github.com/mentat-analytics/backend/issues/2)
