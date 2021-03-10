### Instructions for remote connection to Cloud SQL

- [Configure access to your Cloud SQL instance (add your IP)](https://cloud.google.com/sql/docs/mysql/connect-admin-ip#:~:text=Go%20to%20the%20Cloud%20SQL%20Instances%20page%20in%20the%20Google%20Cloud%20Console.&text=Click%20the%20instance%20to%20open,where%20the%20client%20is%20installed.)
- In [`/.env`](.env) file add `MYSQL_CONN_STR` with value in the format: `root:<password>@tcp(<remote db ip address>:3306)/core?parseTime=true`  
  Replace `<password>` and `<remote db ip address>` with actual values.

### Stripe test card token: "tok_1IT0kfDvMUJnI35tdU1YhOLy"
