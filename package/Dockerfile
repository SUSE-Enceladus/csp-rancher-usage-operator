FROM registry.suse.com/bci/bci-micro:15.4

ARG user=usage-operator

RUN echo "$user:x:1000:1000::/home/$user:/bin/bash" >> /etc/passwd && \
    echo "$user:x:1000:" >> /etc/group && \
    mkdir /home/$user && \
    chown -R $user:$user /home/$user

COPY bin/csp-rancher-usage-operator /usr/bin/

USER $user

CMD ["csp-rancher-usage-operator"]
