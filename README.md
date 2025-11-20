1. ## Getting Started

2. ```git clone https://github.com/shravanasati/scopex-go-assignment```

3. ```cd scopex-go-assignment```

	> [OPTIONAL] If you want attendance report notifications via email, add `RESEND_API_KEY` under environment key of app service in the `docker-compose.yml` file. API key can be obtained from [resend.com](https:///resend.com). 

4. ```docker compose up --build```

5. Run migrations (only need to do this once):

```sh
docker compose exec -T db \
	sh -c "mysql -uhomestead -p!Secret1234 scopex-assignment" < migration.sql
```

Make sure the credentials match those in [properties-prod.yml](./resource/properties-prod.yaml) file (it should already).

6.  Browse Swagger UI [http://localhost:8999/swagger/index.html](http://localhost:8999/swagger/index.html).

	- Login using these credentials on `/api/login`

		* Username: `admin`
		* Password: `admin1234`

	- Copy the response access token and authorize on the UI with the following as the `Authorization` header value: `Bearer <access_token>`.

	- Now you can access all API routes. Try creating a student using the `POST /students` route. Mark their attendance using `POST /attendance/mark` route, get it using `GET /attendance/{student_id}`. Once students and their attendance are created, you'll see attendance reports printed on console and sent on emails if the `RESEND_API_KEY` env var is setup.

7. Run Tests
```
go test -v
```

### Notes

##### Report generation

- Report generation is configured to run every minute for weekly report, and every 2 minutes for monthly report for tesitng purposes. The actual weekly and monthly cron expressions are commented out in [./cronjob/cron_job.go](./cronjob/cron_job.go) file.

- Email notifications are sent only if the `RESEND_API_KEY` is configured.

- Email includes a pretty HTML document that has the student's attendance stats. Each email is sent in background using a goroutine. Synchronization is handled using waitgroups.

##### Optimization

- Report generation avoids N+1 queries using `JOIN`s and `GROUP BY`.

- All DB operations have a time limit of 5 seconds (except report generation, which has 10s timeout), enforced using contexts.

- All listing endpoints using pagination.

- DB indices are created on these columns: `students.id`, `students.email` and a composite index on `attendance.student_id,date`.

- Prepared statements are used everywhere.

- The gin router and DB connection is reused (singleton).

##### Bonus points

- The application is fully dockerized using a multi-stage dockerfile (image size ~52MB, application binary size 41MB).

- Swagger documentation is automatically created using the `swag init` command.

- CI is setup using GitHub Actions which runs test and builds.

- The database migration script `migrate.sh` can be used to apply local migrations.