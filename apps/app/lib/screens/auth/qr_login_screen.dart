import 'package:flutter/material.dart';
import 'package:qr_code_scanner/qr_code_scanner.dart';
import 'package:provider/provider.dart';
import 'package:go_router/go_router.dart';
import 'package:smg_app/providers/auth_provider.dart';
import 'dart:convert';

class QRLoginScreen extends StatefulWidget {
  const QRLoginScreen({super.key});

  @override
  State<QRLoginScreen> createState() => _QRLoginScreenState();
}

class _QRLoginScreenState extends State<QRLoginScreen> {
  final GlobalKey qrKey = GlobalKey(debugLabel: 'QR');
  QRViewController? controller;
  bool isScanning = true;

  @override
  void dispose() {
    controller?.dispose();
    super.dispose();
  }

  void _onQRViewCreated(QRViewController controller) {
    this.controller = controller;
    controller.scannedDataStream.listen((scanData) {
      if (scanData.code != null && isScanning) {
        _handleQRCodeScanned(scanData.code!);
      }
    });
  }

  void _handleQRCodeScanned(String qrCode) async {
    setState(() {
      isScanning = false;
    });

    try {
      final qrData = jsonDecode(qrCode);
      final token = qrData['token'] as String?;
      final expires = qrData['expires'] as int?;

      if (token == null || expires == null) {
        _showError('無效的 QR Code');
        return;
      }

      // Check if QR code has expired
      final now = DateTime.now().millisecondsSinceEpoch ~/ 1000;
      if (now > expires) {
        _showError('QR Code 已過期');
        return;
      }

      final authProvider = Provider.of<AuthProvider>(context, listen: false);
      final success = await authProvider.verifyQRCode(token);

      if (success && mounted) {
        context.go('/home');
      } else if (mounted) {
        _showError('QR Code 驗證失敗');
      }
    } catch (e) {
      _showError('QR Code 格式錯誤');
    }
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.red,
      ),
    );
    
    // Allow scanning again after error
    setState(() {
      isScanning = true;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('QR Code 登入'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: Column(
        children: [
          // Instructions
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Icon(
                          Icons.info,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                        const SizedBox(width: 8),
                        const Text(
                          '如何使用 QR Code 登入',
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    const Text('1. 在網頁版後台點擊個人頭像'),
                    const Text('2. 選擇「QR Code」選項'),
                    const Text('3. 使用此相機掃描螢幕上的 QR Code'),
                    const Text('4. 等待驗證完成'),
                  ],
                ),
              ),
            ),
          ),
          
          // QR Scanner
          Expanded(
            child: Container(
              margin: const EdgeInsets.all(16),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(12),
                child: QRView(
                  key: qrKey,
                  onQRViewCreated: _onQRViewCreated,
                  overlay: QrScannerOverlayShape(
                    borderColor: Theme.of(context).colorScheme.primary,
                    borderRadius: 12,
                    borderLength: 30,
                    borderWidth: 10,
                    cutOutSize: 300,
                  ),
                ),
              ),
            ),
          ),
          
          // Status and controls
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                Consumer<AuthProvider>(
                  builder: (context, authProvider, child) {
                    if (authProvider.isLoading) {
                      return const Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          CircularProgressIndicator(),
                          SizedBox(width: 16),
                          Text('驗證中...'),
                        ],
                      );
                    }
                    
                    return Text(
                      isScanning ? '請對準 QR Code' : '處理中...',
                      style: Theme.of(context).textTheme.bodyLarge,
                    );
                  },
                ),
                const SizedBox(height: 16),
                
                // Flash toggle
                ElevatedButton.icon(
                  onPressed: () async {
                    await controller?.toggleFlash();
                  },
                  icon: const Icon(Icons.flash_on),
                  label: const Text('切換閃光燈'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}