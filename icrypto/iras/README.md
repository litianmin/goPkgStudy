RSA理论知识 （去为知笔记那里看吧）
============================

[参考链接](https://www.cnblogs.com/cjm123/p/8243424.html)



RSA基本知识
============================

[参考链接](https://blog.csdn.net/samsho2/article/details/84255382)
[参考链接2](http://www.361way.com/ras-basics/5820.html)



1）RSA_PKCS1_PADDING 填充模式，最常用的模式


要求:
输入 必须 比 RSA 钥模长(modulus) 短至少11个字节, 也就是 RSA_size(rsa) – 11
    如果输入的明文过长，必须切割， 然后填充
          
输出 和modulus一样长

根据这个要求，对于512bit的密钥， block length = 512/8 – 11 = 53 字节

2） RSA_PKCS1_OAEP_PADDING 
RSA_size(rsa) – 41 

3）for RSA_NO_PADDING 不填充
RSA_size(rsa)

一般来说， 我们只用RSA来加密重要的数据，比如AES的key, 128bits = 16

加密的输出，总是等于key length

对同样的数据，用同样的key进行RSA加密， 每次的输出都会不一样； 但是这些加密的结果都能正确的解密
—————
预备知识
I2OSP – Integer-to-Octet-String primitive 大整数转换成字节串
I2OSP (x, xLen)
输入: x 待转换的非负整数
    xLen 结果字节串的可能长度

————
加密原理 RSAEP ((n, e), m)
输入: (n,e) RSA公钥
m 值为0到n-1 之间一个大整数，代表消息
输出: c 值为0到n-1之间的一个大整数，代表密文
假定： RSA公钥(n,e)是有效的
步骤：
1. 如果m不满足 0 2. 让 c = m^e % n (m的e次幂 除以n ,余数为c)
3. 输出 c

解密原理 RSADP (K, c)
输入：   K RSA私钥，K由下面形式：
           一对(n,d)
一个五元组(p, q, dP, dQ, qInv)
          一个可能为空的三元序列(ri, di, ti), i=3,...,u
c 密文
输出：   m 明文

步骤：
1. 如果密文c不满足 0 < c < n-1, 输出 'ciphertext repersentative out of range'
2. 按照如下方法计算m:
a. 如果使用私钥K的第一种形式(n, d), 就让 m = c^d % n (c的d次幂，除以n,余数为m)
b. 如果使用私钥K的第二种像是(p,q, dP, dQ, qInv)和(ri, di, ti),
--------------

----------------
加密 RSAES-PKCS1-V1_5-ENCRYPT ((n, e), M)

输入： (n, e) 接收者的公开钥匙， k表示n所占用的字节长度
     M     要加密的消息， mLen表示消息的长度 mLen ≤ k – 11

输出： C     密文， 占用字节数 也为 k

步骤：
1.长度检查， 如果 mLen > k-11, 输出 “message too long”

2. EME-PKCS1-v1_5 编码
a) 生成一个 伪随机非零串PS ， 长度为 k – mLen – 3， 所以至少为8， 因为 k-mLen>11
b) 将PS， M，以及其他填充串 一起编码为 EM， 长度为 k, 即:
EM = 0×00 || 0×02 || PS || 0×00 || M 

3.RSA 加密
a)将EM转换成一个大证书m
m = OS2IP(EM)

b)对公钥(n,e) 和 大整数 m, 使用RSAEP加密原理，产生一个整数密文c
c = RSAEP((n,e0, m)

c)将整数c转换成长度为k的密文串
C = I2OSP(c, k)

4.输出密文C
—————-
解密 RSAES-PKCS1-V1_5-DECRYPT (K, C)

输入： K 接收者的私钥
     C   已经加密过的密文串，长度为k (与RSA modulus n的长度一样）
输出： M 消息明文， 长度至多为 k-11

步骤：
1. 长度检查：如果密文C的长度不为k字节（或者 如果 k<11), 输出“decryption error"

2. RSA解密
a. 转换密文C为一个大整数c
c = OS2IP(C)
b. 对RSA私钥(n,d)和密文整数c 实施解密， 产生一个 大整数m
m = RSADP((n,d), c)
如果RSADP输出'ciphertext representative out of range'(意味c>=n), 就输出’decryption error”
c. 转换 m 为长度为k的EM串
     EM = I2OSP(m, k)
3. EME-PKCS1-v1_5 解码：将EM分为 非零的PS串 和 消息 M
     EM = 0×00 || 0×02 || PS || 0×00 || M
如果EM不是上面给出的格式，或者PS的长度小于8个字节， 那么就输出’decryption error’

5. 输出明文消息M

——————–
签名 RSASSA-PSS-SIGN (K, M)
输入   K 签名者的RSA私钥
     M 代签名的消息，一个字节串
输出   S 签名，长度为k的字节串，k是RSA modulus n的字节长度

步骤:
1. EMSA-PSS encoding: 对消息M实施EMSA-PSS编码操作，产生一个长度为 [(modBits -1)/8]的编码消息EM。 整数 OS2IP(EM)的位长最多是 modBits-1, modBits是RSA modulus n的位长度
EM = EMSA-PSS-ENCODE (M, modBits – 1) 

注意：如果modBits-1 能被8整除，EM的字节长比k小1；否则EM字节长为k

2. RSA签名:
a. 将编码后的消息 EM 转换成一个大整数m
m = OS2IP(EM)
b. 对私钥K和消息m 实施 RSASP1 签名，产生一个 大整数s表示的签名
   s = RSASP1 (K, m)
c. 把大整数s转换成 长度为k的字串签名S
S = I2OSP(s, k)
3.输出签名S
———–
验证签名 RSASSA-PSS-VERIFY ((n, e), M, S)
输入： (n,e) 签名者的公钥
M 签名者 发来的消息，一个字串
    S 待验证的签名， 一个长度为k的字串。k是RSA Modulus n的长度
输出： ’valid signature’ 或者 ‘invalid signature’
步骤：
1. 长度检查: 如果签名S的长度不是k, 输出’invalid signature’

2. RSA验证
a) 将签名S转换成一个大整数s
s = OS2IP (S)
b) 对公钥 (n,e) 和 s 实施 RSAVP1 验证， 产生一个 大整数m
m = RSAVP1 ((n, e), s)
c) 将m 转换成编码的消息EM,长度 emLen = [ (modBits -1)/8 ] 字节。 modBits是RSA modulus n的位长
   EM = I2OSP (m, emLen)
注意： 如果 modBits-1可以被8整除，那么emLen = k-1,否则 emLen = k

3. EMSA-PSS验证： 对消息M和编码的EM实施一个 EMSA-PSS验证操作，决定他们是否一致：
Result = EMSA-PSS-VERIFY (M, EM, modBits – 1) 

4. 如果Result = “consistent“，那么输出 ”valid signature”
否则， 输出 ”invalid signature”

———–
签名，还可以使用 EMSA-PKCS1-v1_5 encoding编码方法 来产生 EM:
EM = EMSA-PKCS1-V1_5-ENCODE (M, k)

验证签名是，使用 EMSA-PKCS1-v1_5对 M产生第2个编码消息EM’
   EM’ = EMSA-PKCS1-V1_5-ENCODE (M, k) .

然后比较 EM和EM’ 是否相同

———————

RSA的加密机制有两种方案一个是RSAES-OAEP，另一个RSAES-PKCS1-v1_5。PKCS#1推荐在新的应用中使用RSAES- OAEP，保留RSAES-PKCS#1-v1_5跟老的应用兼容。它们两的区别仅仅在于加密前编码的方式不同。而加密前的编码是为了提供了抵抗各种活动的敌对攻击的安全机制。

PKCS#1的签名机制也有种方案：RSASSA-PSS和RSASSA-PKCS1-v1_5。同样，推荐RSASSA-PSS用于新的应用而RSASSA-PKCS1-v1_5用于兼容老的应用。

——————–

RSAES-OAEP-ENCRYPT ((n, e), M, L)
选项：   Hash 散列函数(hLen 表示 散列函数的输出的字节串的长度）
      MGF 掩码生成函数
输入: (n,e) 接收者的RSA公钥(k表示RSA modulus n的字节长度）
   M 待加密的消息，一个长度为mLen的字节串 mLen <= k - 2 hLen -2
L 同消息关联的可选的标签，如果不提供L，就采用空串
输出: C 密文，字节长度为k
步骤：
1.长度检查
a. 如果L的长度超过 hash函数的输入限制(对于SHA-1, 是2^61 -1），输出 label too long
b. mLen > k – 2hLen -2, 输出 message too long
2. EME-OAEP编码