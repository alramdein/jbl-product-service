
# Setup
Copy `.env.example` and rename it to `.env`. Make sure all the credentials are correct. 
> Please note that if you want to run it on Docker, make sure the `DB_HOST` is the 
> same as in the docker-compose. In this case it must be `'db'`
# How To Run
1. To run in docker:
    ```
    make build
    ```
2. To run locally:
    ```
    make run
    ```
3. To generate mock test:
    ```
    make mockgen
    ```
4. To run test:
    ```
    make test
    ```
3. To run migration up:
    ```
    make migrate-up
    ```
5. To run migration down:
    ```
    make migrate-down
    ```
6. To run migration seed:
    ```
    make seed
    ```
8. To generate swagger docs
    ```
    make swag
    ```