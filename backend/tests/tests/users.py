import unittest
from testutils import IntegrationTestCase

class UsersTestCase(IntegrationTestCase):
    def test_user_signup(self):
        r = self._client.post('/users/signup', json={'email': 'user-signup@example.com', 'password': 'password123'})
        self.assertEqual(r.status_code, 200)

    def test_duplicate_email(self):
        r = self._client.post('/users/signup', json={'email': 'user-signup-duplicate@example.com', 'password': 'password123'})
        self.assertEqual(r.status_code, 200)

        r = self._client.post('/users/signup', json={'email': 'user-signup-duplicate@example.com', 'password': 'password123'})
        self.assertEqual(r.status_code, 400)