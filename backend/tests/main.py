#!/usr/bin/env python3

import argparse
import unittest
import time

import testutils
import tests

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('-p', '--port', type=int)
    parser.add_argument('--host', type=str)
    args = parser.parse_args()

    if args.host is None or args.port is None:
        print('Host or port option are required')
        exit(1)

    client = testutils.Client(args.host, int(args.port))

    wait_for_server_to_be_up(client)

    testloader = testutils.IntegrationTestLoader(client)
    runner = unittest.TextTestRunner()
    result = runner.run(suite(testloader))
    if not result.wasSuccessful():
        exit(1)



def suite(loader: testutils.IntegrationTestLoader) -> unittest.TestSuite:
    suite = unittest.TestSuite()
    suite.addTest(loader.loadTestsFromTestCase(tests.UsersTestCase))
    return suite


def wait_for_server_to_be_up(client: testutils.Client):
    for i in range(30):
        try:
            r = client.get('/health')
            if r.status_code == 200:
                return
        except:
            pass
        time.sleep(1)

    raise 'Server is unreachacble'

if __name__ == "__main__":
    main()