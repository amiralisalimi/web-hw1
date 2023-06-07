from locust import HttpUser, task

class HelloWorldUser(HttpUser):
    @task
    def auth(self):
        self.client.get('localhost:5052/auth')