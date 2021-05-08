# SSH Auto Connector for EC2 Instance

Simple commands which connects to AWS EC2 Instances! 

You don't need to check your changed EC2 instance IP anymore. Just use **'ec2-connect'**

<p align="center"><img src="https://github.com/alicek106/go-ec2-ssh-autoconnect/blob/master/gif/pic.gif" width="70%"></p>

[ec2-ssh-autoconnect](https://github.com/alicek106/ec2-ssh-autoconnect) 의 속도를 개선한 Go 구현 버전입니다. 

# 1. Features

Features that this script provide is..

- Start and stop a EC2 instance
- Start and stop multiple EC2 instances as a group
- **Connect SSH to a EC2 instance automatically**
- List all EC2 instances

# 2. Install

## 2.1 Install binary

Download release binary. Currently, MacOS is only supported.

```
$ wget https://github.com/alicek106/go-ec2-ssh-autoconnect/releases/download/v0.7/ec2-connect-darwin && \
chmod +x ec2-connect-darwin && \
mv ec2-connect-darwin /usr/local/bin/ec2-connect
```

## 2.2 Create configuration file

Create configuration file as **/etc/ec2_connect_config.yaml** like below.

```yaml
version: v1
spec:
  # .spec.credentials is Optional. Instead you can define keys in ~/.aws/credentials
  credentials: 
    accessKey: "... :D"
    secretKey: "... :D"

  privateKeys:
  - name: "default" # You should define default key
    path: "/Users/alicek106/dev/keys/DockerEngineTestInstance.pem"
  - name: "customKey1" # Optional. you can use this key as alternative
    path: "/Users/alicek106/dev/keys/VaultAdminKey.pem"

  instanceGroups:
  - name: "kubeadm_part"
    instances: ["kubeadm_master", "kubeadm_worker0"]
  - name: "kubeadm"
    instances: ["kubeadm_master", "kubeadm_worker0", "kubeadm_worker1", "kubeadm_worker2"]
  - name: "consul"
    instances: ["consul-1", "consul-2", "consul-3"]
  - name: "vault"
    instances: ["vault-1-active", "vault-2-standby"]
```

Configuration file consists of three part.

- **spec.credentails.accessKey and secretKey** : Optional. AWS credentials to use. You can define in ~/.aws/credential as a static credential. Or you can also use assume role from ~/.aws/config (by source_profile and role_arn) to ~/.aws/credentials by shared config.
- **spec.privateKeys** : Default SSH private key when using **ec2-connect connect** for a EC2 instance. By default, privateKey named "default" will be used, but you can define another SSH private key to connect SSH. It is used by **--key** parameter in command
- **spec.instanceGroups** : It defines group for starting and stoping multiple EC2 instances. Above example defined 'kubeadm_part' group, so if you use command **ec2-connect group start kubeadm_part**, it will start two instances (kubeadm-worker, kubeadm-worker0).


## 2.3 Setting AWS credential

AWS credentials priority is like below.

1. Define credentials in bash environment values: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY  (High Priority)
2. Define credentials in /etc/ec2_connect_config.yaml (Medium Priority)
3. Define credentials in ~/.aws/credentials (Low Priority)

FYI, I'm using third way (~/.aws/credentials) :D

# 3. How to use (Easy!)

1. Check EC2 instance list using **ec2-connect list**

   ```
   $ ec2-connect list
   2019/08/11 22:26:08 Cannot found credential in environment variable.
   2019/08/11 22:26:08 Found credential in configuration file.
   2019/08/11 22:26:09 Succeed to validate AWS credential.
   Instance ID		Instance Name		IP Address	Status
   i-04f9e3da95a25e939	kubeadm-master		Unknown		stopped
   i-0c35c716a6442a72e	kubeadm-worker0		Unknown		stopped
   i-0ef65f12957a67559	kubeadm-worker2		Unknown		stopped
   i-0c4bb55ba07c9edbf	kubeadm-worker1		Unknown		stopped
   i-0994dac6654fd59e1	Test			13.209.67.72	running
   ```

2. Start EC2 instance using **ec2-connect start**
   ```
   $ ec2-connect start Test
   2019/08/11 22:26:22 Cannot found credential in environment variable.
   2019/08/11 22:26:22 Found credential in configuration file.
   2019/08/11 22:26:22 Succeed to validate AWS credential.
   2019/08/11 22:26:23 Starting EC2 instance : Test (instance ID: i-0994dac6654fd59e1)
   2019/08/11 22:26:23 Succeed to start EC2 instances.
   ```
3. Connect to EC2 instance by **ec2-connect connect [EC2 instance name]**. This command uses private key defined in /etc/ec2_connect_config.yaml (spec.privateKeys)

   ```
   $ ec2-connect connect Test
   2019/08/11 22:26:54 Cannot found credential in environment variable.
   2019/08/11 22:26:54 Found credential in configuration file.
   2019/08/11 22:26:54 Succeed to validate AWS credential.
   2019/08/11 22:26:55 Instance in active.
   Welcome to Ubuntu 16.04.5 LTS (GNU/Linux 4.4.0-1088-aws x86_64)
   
   ...
   
   New release '18.04.2 LTS' available.
   Run 'do-release-upgrade' to upgrade to it.
   
   
   Last login: Sun Aug 11 13:15:12 2019 from 1.222.77.99
   ubuntu@testbed:~$
   ```
   
   > **Tip** : If a instance is in STOP state, 'connect' command automatically start the instance and connect SSH. So you don't need to command 'start' actually. Just use **connect**!
   
   Or, you can use user-defined key in ec2_connect_config.yaml by specifying --key. By default, this script uses "default" key in config file.

   ```
   $ ec2-connect connect Test --key=MY_CUSTOM_KEY_PATH
   ```


4. Stop EC2 instance by **ec2-connect stop [EC2 instance name]**

   ```
   $ ec2-connect stop Test
   2019/08/11 22:28:12 Cannot found credential in environment variable.
   2019/08/11 22:28:12 Found credential in configuration file.
   2019/08/11 22:28:12 Succeed to validate AWS credential.
   2019/08/11 22:28:12 Stoping EC2 instance : Test (instance ID: i-0994dac6654fd59e1)
   2019/08/11 22:28:13 Succeed to stop EC2 instances.
   ```

5. If you defined **custom group** in /etc/ec2_connect_config.yaml, you can use 'group start' or 'group stop'

   ```
   $ ec2-connect group start kubeadm_part
   2019/08/11 22:29:07 Cannot found credential in environment variable.
   2019/08/11 22:29:07 Found credential in configuration file.
   2019/08/11 22:29:08 Succeed to validate AWS credential.
   2019/08/11 22:29:08 Starting EC2 instance : kubeadm-master (instance ID: i-04f9e3d...)
   2019/08/11 22:29:08 Starting EC2 instance : kubeadm-worker0 (instance ID: i-0c35c71...)
   2019/08/11 22:29:08 Succeed to start EC2 instances.
   ...
   ```
   
