from locust import HttpUser, task

class HelloWorldUser(HttpUser):
    @task
    def auth(self):
        self.client.get('/auth')
    @task
    def get_users(self):
        self.client.post(
            '/get-users',
            json={
                'authKey': '31',
                'userId': '0',
                'withSqlInject': False,
            }
        )