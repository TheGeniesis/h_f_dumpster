from types import ClassMethodDescriptorType
import pexpect
import tempfile
import base64
import sys

def ssh(host, cmd, user, password, timeout=30, bg_run=False):
    """SSH'es to a host using the supplied credentials and executes a command.
    Throws an exception if the command doesn't return 0.
    bgrun: run command in the background"""

    try:
        fname = tempfile.mktemp()
        fout = open(fname, 'w')

        options = '-q -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -oPubkeyAuthentication=no -t'
        if bg_run:
            options += ' -f'
        ssh_cmd = 'ssh %s@%s %s "%s"' % (user, host, options, cmd)
        child = pexpect.spawnu(ssh_cmd, timeout=timeout)
        child.expect(['[pP]assword:'])
        child.sendline(password)
        child.logfile = fout
        child.expect(pexpect.EOF)
        child.close()
        fout.close()

        fin = open(fname, 'r')
        stdout = fin.read()
        fin.close()

        if 0 != child.exitstatus:
            print("problem with host " + host)
            #raise Exception(stdout)

            return
    except:
        print("problem with host " + host)
        return

    print("Executed succesfully for  " + host)

    return stdout

# Update these lines
host_list = ['']
user = ''

f = open(sys.path[0]+'/password.txt', 'r')
password = base64.b64decode(f.read()).decode('utf-8')

for host in host_list:
    ssh(host=host, cmd='sudo touch test12345.txt; rm -f test12345.txt', user=user, password=password)
