import requests

class Client:
    def __init__(self, host: str, port: int):
        self._host = host
        self._port = port

    def post(self, endpoint: str, json=None, headers=[], params=[]) -> requests.Response:
        return requests.post(self._format_url(endpoint), json=json)

    def get(self, endpoint: str, params=[], headers=[]) -> requests.Response:
        return requests.get(self._format_url(endpoint), params=params, headers=headers)

    def _format_url(self, endpoint: str) -> str:
        if endpoint.startswith('/'):
            endpoint = endpoint[1:]
        return f'http://{self._host}:{self._port}/{endpoint}'