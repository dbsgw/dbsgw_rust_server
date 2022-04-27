create table user_auth
(
    id            bigint auto_increment
        primary key,
    uid           varchar(22)      default '0' not null comment '用户id',
    identity_type tinyint unsigned default '1' not null comment '1邮箱 2gitee 3githup ',
    identifier    varchar(50)      default ''  not null comment '手机号 邮箱 用户名或第三方应用的唯一标识',
    certificate   varchar(20)      default ''  not null comment '密码凭证(站内的保存密码，站外的不保存或保存token)',
    create_time   int unsigned     default '0' not null comment '绑定时间',
    update_time   int unsigned     default '0' not null comment '更新绑定时间',
    constraint only
        unique (uid, identity_type)
)
    comment '用户授权表';

create index idx_uid
    on user_auth (uid);

create table user_base
(
    uid              varchar(22)                  not null comment '用户ID'
        primary key,
    user_role        tinyint unsigned default '2' not null comment '2正常用户 3禁言用户 4虚拟用户 5运营',
    register_source  tinyint unsigned default '0' not null comment '注册来源：1邮箱 2gitee 3githup ',
    user_name        varchar(32)      default ''  not null comment '用户账号，必须唯一',
    nick_name        varchar(32)      default ''  not null comment '用户昵称',
    gender           tinyint unsigned default '0' not null comment '用户性别 0-female 1-male',
    birthday         bigint unsigned  default '0' not null comment '用户生日',
    signature        varchar(255)     default ''  not null comment '用户个人签名',
    mobile           varchar(16)      default ''  not null comment '手机号码(唯一)',
    mobile_bind_time int unsigned     default '0' not null comment '手机号码绑定时间',
    email            varchar(100)     default ''  not null comment '邮箱(唯一)',
    email_bind_time  int unsigned     default '0' not null comment '邮箱绑定时间',
    face             varchar(255)     default ''  not null comment '头像',
    face200          varchar(255)     default ''  not null comment '头像 200x200x80',
    create_time      int unsigned                 not null comment '创建时间',
    srcface          varchar(255)     default ''  not null comment '原图头像',
    update_time      int unsigned                 not null comment '修改时间',
    constraint user_base_uid_uindex
        unique (uid)
)
    comment '用户基础信息表';

create table user_location
(
    uid           varchar(22)              not null comment '用户ID'
        primary key,
    curr_nation   varchar(10)  default ''  not null comment '所在地国',
    curr_province varchar(10)  default ''  not null comment '所在地省',
    curr_city     varchar(10)  default ''  not null comment '所在地市',
    curr_district varchar(20)  default ''  not null comment '所在地地区',
    location      varchar(255) default ''  not null comment '具体地址',
    longitude     decimal(10, 6)           null comment '经度',
    latitude      decimal(10, 6)           null comment '纬度',
    update_time   int unsigned default '0' null comment '修改时间'
)
    comment '用户定位表';


