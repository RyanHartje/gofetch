// thanks to: https://radu-matei.com/blog/kubernetes-jenkins-azure/
podTemplate(label: 'mypod', containers: [
    containerTemplate(name: 'docker', image: 'docker', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'kubectl', image: 'lachlanevenson/k8s-kubectl:v1.8.0', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'helm', image: 'lachlanevenson/k8s-helm:latest', command: 'cat', ttyEnabled: true)
  ],
  volumes: [
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
  ]) {
    node('mypod') {

        stage('do some Docker work') {
            container('docker') {

                    sh """
                        docker pull golang:1.9-alpine
                       """
                }
            }
        }

        stage('do some kubectl work') {
            container('kubectl') {
                    sh "kubectl get nodes"
                }
            }
        }
        stage('do some helm work') {
            container('helm') {
               sh "helm ls"
            }
        }
    }
}
