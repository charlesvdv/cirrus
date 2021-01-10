import unittest
from testutils import IntegrationTestCase

class UserSessionTestCase(IntegrationTestCase):
    def test_session_authentication(self):
        r = self._client.post('/users', json={'email': 'user-auth@example.com', 'password': 'password123'})
        self.assertEqual(r.status_code, 200)

        r = self._client.post('/session/authenticate', json={'email': 'user-auth@example.com', 'password': 'password123'})
        self.assertEqual(r.status_code, 200)
        self.assertNotEqual(r.json()["access_token"]["token"], "")
        self.assertNotEqual(r.json()["refresh_token"], "")
        self.assertNotEqual(r.json()["client_reference"], "")

        access_token = r.json()["access_token"]["token"]
        r = self._client.get('/users', headers={'Authorization': f'Bearer {access_token}'})
        self.assertEqual(r.status_code, 200)
        self.assertEqual(r.json()["email"], 'user-auth@example.com')

    def test_unknown_session(self):
        r = self._client.get('/users', headers={'Authorization': 'Bearer unknown-random-token'})
        self.assertEqual(r.status_code, 401)

    def test_session_authentication_with_invalid_password(self):
        r = self._client.post('/users', json={'email': 'user-auth-invalid-password@example.com', 'password': 'password123'})
        self.assertEqual(r.status_code, 200)

        r = self._client.post('/session/authenticate', json={'email': 'user-auth-invalid-password@example.com', 'password': 'invalidpassword'})
        self.assertEqual(r.status_code, 400)
        self.assertEqual(r.json()["message"], "Invalid user name or password")

    def test_session_authentication_with_invalid_email(self):
        r = self._client.post('/users', json={'email': 'user-auth-invalid-email@example.com', 'password': 'password123'})
        self.assertEqual(r.status_code, 200)

        r = self._client.post('/session/authenticate', json={'email': 'random-invalid-email@example.com', 'password': 'password123'})
        self.assertEqual(r.status_code, 400)
        self.assertEqual(r.json()["message"], "Invalid user name or password")

