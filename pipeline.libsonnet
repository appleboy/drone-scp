{
  test:: {
    kind: 'pipeline',
    name: 'testing',
    platform: {
      os: 'linux',
      arch: 'amd64',
    },
    steps: [
      {
        name: 'vet',
        image: 'golang:1.18',
        pull: 'always',
        commands: [
          'make vet',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      // {
      //   name: 'test',
      //   image: 'golang:1.18-alpine',
      //   pull: 'always',
      //   commands: [
      //     'apk add git make curl perl bash build-base zlib-dev ucl-dev',
      //     'make ssh-server',
      //     'make test',
      //     'make coverage',
      //   ],
      //   volumes: [
      //     {
      //       name: 'gopath',
      //       path: '/go',
      //     },
      //   ],
      // },
      {
        name: 'codecov',
        image: 'robertstettner/drone-codecov',
        pull: 'always',
        settings: {
          token: { 'from_secret': 'codecov_token' },
        },
      },
    ],
    volumes: [
      {
        name: 'gopath',
        temp: {},
      },
    ],
  },

  build(name, os='linux', arch='amd64'):: {
    kind: 'pipeline',
    name: os + '-' + arch,
    platform: {
      os: os,
      arch: arch,
    },
    steps: [
      {
        name: 'build-push',
        image: 'golang:1.18',
        pull: 'always',
        environment: {
          CGO_ENABLED: '0',
        },
        commands: [
          'go build -v -ldflags \'-X main.build=${DRONE_BUILD_NUMBER}\' -a -o release/' + os + '/' + arch + '/' + name,
        ],
        when: {
          event: {
            exclude: [ 'tag' ],
          },
        },
      },
      {
        name: 'build-tag',
        image: 'golang:1.18',
        pull: 'always',
        environment: {
          CGO_ENABLED: '0',
        },
        commands: [
          'go build -v -ldflags \'-X main.version=${DRONE_TAG##v} -X main.build=${DRONE_BUILD_NUMBER}\' -a -o release/' + os + '/' + arch + '/' + name,
        ],
        when: {
          event: [ 'tag' ],
        },
      },
      {
        name: 'executable',
        image: 'golang:1.18',
        pull: 'always',
        commands: [
          './release/' + os + '/' + arch + '/' + name + ' --help',
        ],
      },
      {
        name: 'dryrun',
        image: 'plugins/docker:' + os + '-' + arch,
        pull: 'always',
        settings: {
          daemon_off: false,
          dry_run: true,
          tags: os + '-' + arch,
          dockerfile: 'docker/Dockerfile.' + os + '.' + arch,
          repo: 'appleboy/' + name,
          cache_from: 'appleboy/' + name,
        },
        when: {
          event: [ 'pull_request' ],
        },
      },
      {
        name: 'publish',
        image: 'plugins/docker:' + os + '-' + arch,
        pull: 'always',
        settings: {
          daemon_off: 'false',
          auto_tag: true,
          auto_tag_suffix: os + '-' + arch,
          dockerfile: 'docker/Dockerfile.' + os + '.' + arch,
          repo: 'appleboy/' + name,
          cache_from: 'appleboy/' + name,
          username: { 'from_secret': 'docker_username' },
          password: { 'from_secret': 'docker_password' },
        },
        when: {
          event: {
            exclude: [ 'pull_request' ],
          },
        },
      },
    ],
    depends_on: [
      'testing',
    ],
    trigger: {
      ref: [
        'refs/heads/master',
        'refs/pull/**',
        'refs/tags/**',
      ],
    },
  },

  release:: {
    kind: 'pipeline',
    name: 'release-binary',
    platform: {
      os: 'linux',
      arch: 'amd64',
    },
    steps: [
      {
        name: 'build-all-binary',
        image: 'golang:1.18',
        pull: 'always',
        commands: [
          'make release'
        ],
        when: {
          event: [ 'tag' ],
        },
      },
      {
        name: 'deploy-all-binary',
        image: 'plugins/github-release',
        pull: 'always',
        settings: {
          files: [ 'dist/release/*' ],
          api_key: { 'from_secret': 'github_release_api_key' },
        },
        when: {
          event: [ 'tag' ],
        },
      },
    ],
    depends_on: [
      'testing',
    ],
    trigger: {
      ref: [
        'refs/tags/**',
      ],
    },
  },

  notifications(os='linux', arch='amd64', depends_on=[]):: {
    kind: 'pipeline',
    name: 'notifications',
    platform: {
      os: os,
      arch: arch,
    },
    steps: [
      {
        name: 'manifest',
        image: 'plugins/manifest',
        pull: 'always',
        settings: {
          username: { from_secret: 'docker_username' },
          password: { from_secret: 'docker_password' },
          spec: 'docker/manifest.tmpl',
          ignore_missing: true,
        },
      },
    ],
    depends_on: depends_on,
    trigger: {
      ref: [
        'refs/heads/master',
        'refs/tags/**',
      ],
    },
  },

  signature(key):: {
    kind: 'signature',
    hmac: key,
  }
}
