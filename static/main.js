const Login = {
    template: '<div>failed to load login page</div>',
    data: function() {
        return {
            loginParam: {
                username: "",
                password: ""
            },
            loginLoading: false
        }
    },
    methods: {
        login() {
            if (!this.loginParam.username || !this.loginParam.password)
                return
            this.loginLoading = true
            this.$http.post("/login.json", this.loginParam).then((response) => {
                this.$store.commit("updateUserInfo", response.data)
                let expireDays = 1000 * 60 * 60 * 24 * 15;
                var exdate = new Date();
                exdate.setDate(exdate.getDate + expireDays);
                document.cookie = "token" + "=" + escape(response.data.token) +
                    ((expireDays == null) ? "" : ";expires=" + exdate.toGMTString()) +
                    ";path=/";
                this.$router.push({ path: '/' })
            }, (response) => {
                this.loginLoading = false
                alert(response.data.ErrMessage)
            })
        }
    }
}


const Home = {
    template: '<div>failed to load home page</div>',
    name: 'home',
    data: function() {
        return {}
    },
    methods: {
        logout() {
            this.delCookie('token');
            this.$router.push('/login');
        },
        toCreateProject() {
            this.$router.push('/createproject');
        },
        toHome() {
            this.$router.push('/');
        }
    }
}

const CreateProject = {
    template: '<div>failed to load create project page</div>',
    name: 'createproject',
    data: function() {
        return {
            configProjectParam: {
                projectName: '',
                gitURL: '',
                httpPort: '',
                projectDesc: '',
                buildCmd: '',
                targetZip: '',
                unzipDir: '',
                programLanguage: '',
                buildDependency: '',
                startCmd: '',
                stopCmd: '',
                preCmd: ''
            },
            buttonStr: '创建'
        }
    },
    methods: {
        back() {
            this.$router.back();
        },
        update() {
            if (!this.configProjectParam.projectName ||
                !this.configProjectParam.gitURL ||
                !this.configProjectParam.buildCmd ||
                !this.configProjectParam.targetZip ||
                !this.configProjectParam.unzipDir ||
                !this.configProjectParam.programLanguage ||
                !this.configProjectParam.startCmd ||
                !this.configProjectParam.stopCmd
            ) {
                alert("请填入必填字段")
                return
            }
            this.$http.post("/create_project.json", this.configProjectParam).then((response) => {
                // 将新创建的项目添加到左侧边栏
                userInfo = this.$store.state.userInfo
                userInfo.projects.push(this.configProjectParam.projectName)
                userInfo.projects.sort()
                this.$store.commit("updateUserInfo", userInfo)
                alert("创建成功");
                // 自动跳转到新项目页
                this.$router.push("/projects/" + this.configProjectParam.projectName);
            }, (response) => {
                alert(response.data.ErrMessage)
            })
        }
    }
}

const UpdateProject = {
    template: '<div>failed to load update project page</div>',
    name: 'updateproject',
    data: function() {
        return {
            configProjectParam: {
                projectName: this.$route.params.project,
                gitURL: '',
                httpPort: '',
                projectDesc: '',
                buildCmd: '',
                targetZip: '',
                unzipDir: '',
                programLanguage: '',
                buildDependency: '',
                startCmd: '',
                stopCmd: '',
                preCmd: ''
            },
            buttonStr: '修改'
        }
    },
    created: function() {
        this.$http.get("/project_config.json", { params: { proj: this.$route.params.project } }).then((response) => {
            this.configProjectParam = response.data
        }, (response) => {
            alert(response.data.ErrMessage)
            this.$router.back();
        })
    },
    methods: {
        back() {
            this.$router.back();
        },
        update() {
            if (!this.configProjectParam.projectName ||
                !this.configProjectParam.gitURL ||
                !this.configProjectParam.buildCmd ||
                !this.configProjectParam.targetZip ||
                !this.configProjectParam.unzipDir ||
                !this.configProjectParam.programLanguage ||
                !this.configProjectParam.startCmd ||
                !this.configProjectParam.stopCmd
            ) {
                alert("请填入必填字段")
                return
            }
            this.$http.post("/update_project.json", this.configProjectParam).then((response) => {
                alert("修改成功");
                this.$router.back();
            }, (response) => {
                alert(response.data.ErrMessage)
            })
        }
    }
}

const Project = {
    template: '<div>failed to load project page</div>',
    name: 'project',
    data: function() {
        return {
            runtimeInfo: {
                Regions: [{
                    Name: '',
                    Envs: null
                }]
            },
            createEnvDialogVisible: false,
            createEnvParam: {
                Name: '',
                ProjName: '',
                Region: '',
                CodeBranch: '',
                CodeVersion: '',
                NodeNum: 0,
                CPU: 0,
                Memory: 0,
            },
        }
    },
    created() {
        this.getProjectRuntime()
    },
    watch: {
        "$route": 'getProjectRuntime'
    },
    methods: {
        back() {
            this.$router.push("/");
        },
        toUpdateProject() {
            this.$router.push('/projects/' + this.$route.params.project + '/update');
        },
        toDeploy(region, env) {
            this.$router.push('/projects/' + this.$route.params.project + '/deploy/' + region + '/' + env);
        },
        getProjectRuntime() {
            this.$http.get("/project_runtime.json", { params: { proj: this.$route.params.project } }).then((response) => {
                this.runtimeInfo = response.data
            }, (response) => {
                alert(response.data.ErrMessage)
                this.$router.back();
            })
        },
        createEnv(region) {
            this.createEnvDialogVisible = false
            this.createEnvParam.ProjName = this.$route.params.project
            this.createEnvParam.Region = region
            this.$http.post("/create_env.json", this.createEnvParam).then((response) => {
                // 更新显示
                this.getProjectRuntime()
                alert("创建成功");
            }, (response) => {
                alert(response.data.ErrMessage)
            })
        }
    }
}

const Deploy = {
    template: '<div>failed to load deploy page</div>',
    name: 'image',
    data: function() {
        return {
            createEnvDialogVisible: false,
            envType: '',
            branch: '',
            cpuNum: 1,
            memNum: 1,
            images: [],
            node_num: 0
        }
    },
    created() {
        this.$http.get("/get_env_node_num.json", {
            params: {
                proj: this.$route.params.project,
                env: this.$route.params.env,
                region: this.$route.params.region
            }
        }).then((response) => {
            this.node_num = response.data
        }, (response) => {
            alert("获取节点数失败: " + response.data.ErrMessage)
            this.$router.back();
        })
        this.getImages()
    },
    methods: {
        back() {
            this.$router.back();
        },
        getImages() {
            this.$http.get("/images.json", {
                params: {
                    proj: this.$route.params.project,
                    env: this.$route.params.env
                }
            }).then((response) => {
                this.images = response.data
            }, (response) => {
                alert(response.data.ErrMessage)
                this.$router.back();
            })
        },
        deploy(id) {
            this.$http.post("/update_env_code.json", {
                proj: this.$route.params.project,
                env: this.$route.params.env,
                region: this.$route.params.region,
                image_id: id
            }).then((response) => {
                alert("部署任务已提交");
                this.$router.back();
            }, (response) => {
                alert(response.data.ErrMessage)
            })
        },
        updateNodeNum() {
            this.$http.post("/set_env_node_num.json", {
                proj: this.$route.params.project,
                env: this.$route.params.env,
                region: this.$route.params.region,
                node_num: this.node_num
            }).then((response) => {
                alert("修改任务已提交");
                this.$router.back();
            }, (response) => {
                alert(response.data.ErrMessage)
            })
        }
    }
}


req = new XMLHttpRequest();
req.open('GET', 'login.html', false);
req.send(null);
Login.template = req.responseText;

req = new XMLHttpRequest();
req.open('GET', 'home.html', false);
req.send(null);
Home.template = req.responseText;

req = new XMLHttpRequest();
req.open('GET', 'config-project.html', false);
req.send(null);
CreateProject.template = req.responseText;

// use the same html with CreateProject
UpdateProject.template = req.responseText;

req = new XMLHttpRequest();
req.open('GET', 'project.html', false);
req.send(null);
Project.template = req.responseText;

req = new XMLHttpRequest();
req.open('GET', 'deploy.html', false);
req.send(null);
Deploy.template = req.responseText;


const router = new VueRouter({
    routes: [{
        path: '/login',
        component: Login
    }, {
        path: '/',
        component: Home
    }, {
        path: '/test',
        component: Home
    }, {
        path: '/',
        component: Home,
        children: [{
            path: 'createproject',
            component: CreateProject
        }, {
            path: 'projects/:project/update',
            component: UpdateProject
        }, {
            path: 'projects/:project',
            component: Project
        }, {
            path: 'projects/:project/deploy/:region/:env',
            component: Deploy
        }]
    }, {
        path: '*',
        redirect: '/'
    }]
})

const store = new Vuex.Store({
    state: {
        domain: 'http://127.0.0.1:8000', //保存后台请求的地址，修改时方便（比方说从测试服改成正式服域名）
        userInfo: { //保存用户信息
            name: null,
            token: null,
            groups: null,
            projects: []
        }
    },
    mutations: {
        //更新用户信息
        updateUserInfo(state, newUserInfo) {
            state.userInfo = newUserInfo
            Lockr.set('userinfo', newUserInfo)
        },
        getInfoFromLocker(state) {
            userInfo = Lockr.get('userinfo')
            if (userInfo) {
                state.userInfo = userInfo
            }
        }
    }
})

//获取cookie
getCookie = function(name) {
    var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if (arr = document.cookie.match(reg)) {
        return (arr[2]);
    }
    return null;
};

Vue.prototype.getCookie = getCookie
Vue.prototype.delCookie = (name) => {
    var exp = new Date();
    exp.setTime(exp.getTime() - 1);
    var cval = this.getCookie(name);
    if (cval != null)
        document.cookie = name + "=" + cval + ";expires=" + exp.toGMTString() + ";path=/";
}

new Vue({
    created() {
        this.checkLogin();
    },
    watch: {
        "$route": 'checkLogin'
    },
    el: '#app',
    data: {},
    router,
    store,
    methods: {
        checkLogin() {
            if (!this.getCookie('token')) {
                //如果没有登录状态则跳转到登录页
                this.$router.push('/login');
            } else {
                //如果已经登录则导入本地信息
                this.$store.commit('getInfoFromLocker')
            }
        }
    }
})