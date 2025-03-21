#!/system/bin/sh

PHONE_ID="$1"
BACKUP_DIR="/data/media/0/bak/$PHONE_ID"
TAR_BIN="/data/adb/magisk/busybox tar"
GMS_PKGS=("com.google.android.gms" "com.google.android.gsf" "com.android.vending")
MSG_PKG="com.google.android.apps.messaging"

echo "💾 开始还原 Google 服务 & Messages 应用数据"
echo "📱 当前标识：$PHONE_ID"

# === 环境准备 ===
echo "🔧 暂时关闭 SELinux & 网络"
su -c "setenforce 0"
svc wifi disable
svc data disable

# === 停止 GMS 相关服务 ===
echo "⛔ 停止 GMS & Messages 应用"
for pkg in "${GMS_PKGS[@]}"; do
  if pm list packages | grep -q "$pkg"; then
    am force-stop "$pkg"
    pm disable "$pkg"
    echo "✅ 停止 $pkg"
  fi
done
am force-stop "$MSG_PKG"

# === 拷贝 fake_fingerprint.json（如有）===
if [ -f "$BACKUP_DIR/fake_fingerprint.json" ]; then
  cp "$BACKUP_DIR/fake_fingerprint.json" /data/local/tmp/
  echo "✅ fake_fingerprint.json 已复制"
else
  echo "⚠️ fake_fingerprint.json 不存在，跳过"
fi

# === 清理旧数据（提示性删除）===
echo "🧹 正在清理旧 GMS 数据（保留 APK / JAR / ODEX）..."
for dir in /data/data/com.google.android.gms /data/data/com.google.android.gsf /data/data/com.android.vending; do
  if [ -d "$dir" ]; then
    echo "📂 清理目录: $dir"
    su -c "find $dir -type f ! -name '*.jar' ! -name '*.apk' ! -name '*.odex' -delete"
    su -c "find $dir -type d -empty -delete"
  fi
done

# === 还原 GMS 数据 ===
GMS_BACKUP="$BACKUP_DIR/GoogleBackup.tar.gz"
if [ -f "$GMS_BACKUP" ]; then
  echo "📦 解压 GMS 数据..."
  su -c "$TAR_BIN -xzvf $GMS_BACKUP -C /data/data/"

  for pkg in "${GMS_PKGS[@]}"; do
    uid=$(cmd package list packages -U | grep "$pkg" | grep -oE "uid:[0-9]+" | cut -d: -f2)
    su -c "restorecon -Rv /data/data/$pkg"
    su -c "chown -R $uid:$uid /data/data/$pkg"
    echo "✅ 恢复完成：$pkg"
  done
else
  echo "⚠️ 未找到 GMS 备份文件，跳过..."
fi

# === 还原 Messages 数据 ===
MSG_BACKUP="$BACKUP_DIR/MessagesBackup.tar.gz"
echo "🧼 清空 Messages 缓存..."
pm clear "$MSG_PKG"

if [ -f "$MSG_BACKUP" ]; then
  echo "📦 解压 Messages 数据..."
  su -c "$TAR_BIN -xzvf $MSG_BACKUP -C /data/data/"

  uid=$(cmd package list packages -U | grep "$MSG_PKG" | grep -oE "uid:[0-9]+" | cut -d: -f2)
  su -c "restorecon -Rv /data/data/$MSG_PKG"
  su -c "chown -R $uid:$uid /data/data/$MSG_PKG"

  echo "✅ Messages 数据已恢复"
else
  echo "⚠️ 未找到 Messages 备份文件，跳过..."
fi

# === 恢复 App 与系统 ===
echo "🔁 恢复应用与网络状态"

for pkg in "${GMS_PKGS[@]}"; do
  if pm list packages | grep -q "$pkg"; then
    pm enable "$pkg"
    echo "✅ 启用 $pkg"
  fi
done

svc wifi enable
svc data enable
su -c "setenforce 1"

echo "🎉 Google 服务与消息应用数据恢复完成！"
