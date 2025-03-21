#!/system/bin/sh

# === 基本设置 ===
PHONE_ID="$1"
BACKUP_DIR="/data/media/0/bak/$PHONE_ID"
TAR_BIN="/data/adb/magisk/busybox tar"
EXCLUDE_FILE="/data/local/tmp/tar_exclude.txt"

GMS_PACKAGES=("com.google.android.gms" "com.google.android.gsf" "com.android.vending")
MSG_PACKAGE="com.google.android.apps.messaging"

echo "📦 开始备份 Google 服务 & Messages 应用数据"
echo "📱 当前标识：$PHONE_ID"

# === 准备工作 ===
echo "*.jar" > $EXCLUDE_FILE
echo "*.apk" >> $EXCLUDE_FILE
echo "*.odex" >> $EXCLUDE_FILE

su -c "mkdir -p $BACKUP_DIR"
su -c "chmod 777 $BACKUP_DIR"

# === 临时关闭 SELinux & 网络 ===
echo "🔧 暂时关闭 SELinux 与网络连接"
su -c "setenforce 0"
svc wifi disable
svc data disable

# === 停止 GMS 和 Messages 应用 ===
echo "⛔ 正在关闭相关 App..."
for pkg in "${GMS_PACKAGES[@]}"; do
  if pm list packages | grep -q "$pkg"; then
    am force-stop "$pkg"
    pm disable "$pkg"
    echo "✅ 已关闭 $pkg"
  else
    echo "❌ 未安装 $pkg，跳过..."
  fi
done

am force-stop "$MSG_PACKAGE"

# === 备份 GMS 数据（覆盖式）===
echo "📂 开始备份 GMS 数据 → GoogleBackup.tar.gz"

if [ -d "/data/data/com.google.android.gms" ]; then
  su -c "$TAR_BIN -czvf $BACKUP_DIR/GoogleBackup.tar.gz \
      --numeric-owner --exclude-from=$EXCLUDE_FILE -C /data/data \
      com.google.android.gms com.google.android.gsf com.android.vending"
  echo "✅ GMS 数据备份完成"
else
  echo "⚠️ 未找到 GMS 目录，跳过..."
fi

# === 备份 Messages 数据 ===
echo "📂 开始备份 Messages 数据 → MessagesBackup.tar.gz"

if [ -d "/data/data/$MSG_PACKAGE" ]; then
  su -c "$TAR_BIN -czvf $BACKUP_DIR/MessagesBackup.tar.gz \
      --numeric-owner --exclude-from=$EXCLUDE_FILE -C /data/data \
      $MSG_PACKAGE"
  echo "✅ Messages 数据备份完成"
else
  echo "⚠️ 未找到 Messages 目录，跳过..."
fi

# === 附加文件备份（如 fake_fingerprint.json）===
if [ -f "/data/local/tmp/fake_fingerprint.json" ]; then
  cp /data/local/tmp/fake_fingerprint.json "$BACKUP_DIR/"
  echo "✅ 附加信息 fake_fingerprint.json 已保存"
fi

# === 恢复系统状态 ===
echo "🔁 恢复系统状态..."

for pkg in "${GMS_PACKAGES[@]}"; do
  if pm list packages | grep -q "$pkg"; then
    pm enable "$pkg"
    echo "✅ 已启用 $pkg"
  fi
done

svc wifi enable
svc data enable
su -c "setenforce 1"

rm -f "$EXCLUDE_FILE"

echo "🎉 完成！所有备份已覆盖保存于：$BACKUP_DIR"
