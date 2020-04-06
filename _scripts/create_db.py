#!/usr/bin/python3

# This creates a new DB cluster. It requires the user to input the
# master password via the command line.

import argparse
import os
import subprocess
import sys

def main():
    print('Enter the master database password:')
    password = input()
    if not password or len(password) < 16:
        print('Input must be 16 characters or longer')
        return 1
    
    print('master password: %s', password)
    return subprocess.call([
        'aws',
        '--profile', 'delphis',
        'rds',
        'create-db-cluster',
        '--db-cluster-identifier', 'chatham-staging-aurora-pgsql',
        '--engine', 'aurora-postgresql',
        '--engine-version', '11.6',
        '--master-username', 'postgres',
        '--master-user-password', password,
        '--availability-zones', 'us-west-2a',
        '--database-name', 'chatham',
    ])
    
    return 0

if __name__ == '__main__':
    sys.exit(main())