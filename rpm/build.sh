#!/bin/bash
echo "Recreating rpmbuild directory"
rm -rvf /root/rpmbuild/
rpmdev-setuptree
echo "Building SRPM"
rpmbuild --undefine=_disable_source_fetch -bs /project/gpt-cli-chat/rpm/gpt-cli-chat.spec
mkdir -p ~/.config
mv /project/gpt-cli-chat/copr ~/.config/copr
copr-cli build gpt-cli-chat /root/rpmbuild/SRPMS/*.src.rpm
