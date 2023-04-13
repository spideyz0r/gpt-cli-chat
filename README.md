# gpt-cli-chat [![CI](https://github.com/spideyz0r/gpt-cli-chat/workflows/gotester/badge.svg)][![CI](https://github.com/spideyz0r/gpt-cli-chat/workflows/goreleaser/badge.svg)][![CI](https://github.com/spideyz0r/gpt-cli-chat/workflows/rpm-builder/badge.svg)]
gpt-cli-chat is a cli tool to use open ai chat models.

## Install

### RPM
```
dnf copr enable brandfbb/gpt-cli-chat
dnf install gpt-cli-chat
```

### From source
```
go build -v -o gpt-cli-chat
```

## Usage
```
# gpt-cli-chat
Usage: gpt-cli-chat [-hs] [-a value] [-d value] [-r value] [-t value] [-w value] [parameters ...]
 -a, --api=value  API key (default: OPENAI_API_KEY environment variable)
 -d, --delimiter=value
                  set the delimiter for the user input (default: new line)
 -h, --help       display this help
 -r, --system-role=value
                  system role (default: You're an expert in everything. You
                  like speaking.)
 -s, --stdin      read the message from stdin and exit (default: false)
 -t, --temperature=value
                  temperature (default: 0.8)
 -w, --output-width=value
                  output width (default: 80)
```

## Examples
Execute a command in multiple VMs filtering by their name
```
echo "What do you know about the force?" | gpt-cli-chat -r "You're a jedi" -s
Bot: As an AI language model, I have been programmed with knowledge on the Star
Wars universe, including the Force. The Force is an energy field that binds
together all living things in the galaxy. It can be harnessed by certain
individuals known as Force users, such as Jedi and Sith, to manipulate objects,
influence minds, and even see glimpses of the future. The Force has two aspects:
the light side, which promotes peace and selflessness, and the dark side, which
promotes aggression and selfishness. Jedi are trained to use the Force for good
and to resist the temptations of the dark side.

```

List VMs that match a given filter
```
gpt-cli-chat -t 0.8 -d ';end'
You (press ;end to finish): Hey, would you please,
tell me what does the command kubectl do?;end
Bot: Sure, Kubectl is a command-line tool used to deploy, manage, and monitor
applications in Kubernetes clusters. It allows users to interact with the
Kubernetes API server and perform various operations such as deploying and
scaling applications, inspecting and updating the cluster state, and managing
networking and storage resources. Kubectl is an essential tool for anyone
working with Kubernetes, whether you're a developer, DevOps engineer, or a
system administrator.

```
